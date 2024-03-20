package app

import (
	"sync"
	"time"

	"github.com/aredoff/reagate/internal/config"
	"github.com/aredoff/reagate/internal/log"
	"github.com/aredoff/reagate/pkg/httptracer"
)

func New() *URLMonitor {
	clear_cron := time.Duration(config.Config.GetInt("clear_cron")) * time.Second
	log.Infof("Set clearning cron : %d seconds", int(clear_cron.Seconds()))

	clear_interval := time.Duration(config.Config.GetInt("clear_interval")) * time.Second
	log.Infof("Set clearning interval : %d seconds", int(clear_interval.Seconds()))

	interval := time.Duration(config.Config.GetInt("interval")) * time.Second
	log.Infof("Set monitoring interval: %d seconds", int(interval.Seconds()))

	timeout := time.Duration(config.Config.GetInt("timeout")) * time.Second
	log.Infof("Set monitoring timeout: %d seconds", int(timeout.Seconds()))

	user_agent := config.Config.GetString("user_agent")
	log.Infof("Set user agent: %s", user_agent)

	um := URLMonitor{
		tracerPool: sync.Pool{
			New: func() interface{} {
				tracer := httptracer.New()
				tracer.SetTimeout(timeout)
				tracer.SetHeaders("user-agent", user_agent)
				return tracer
			},
		},
		sites:          make(map[string]*site),
		mu:             &sync.RWMutex{},
		clearTicker:    time.NewTicker(clear_cron),
		interval:       interval,
		clear_interval: clear_interval,
	}

	go um.clearCron()

	return &um
}

type URLMonitor struct {
	tracerPool     sync.Pool
	sites          map[string]*site
	interval       time.Duration
	clear_interval time.Duration
	clearTicker    *time.Ticker
	mu             *sync.RWMutex
}

func (um *URLMonitor) SetInterval(interval time.Duration) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.interval = interval
}

func (um *URLMonitor) Trace(url, method string) *httptracer.TracerResult {
	report, ok := um.getReport(url)
	if ok {
		if report != nil {
			return report
		}
		return um.trace(url, method)
	}
	um.addSite(url)
	return um.trace(url, method)
}

func (e *URLMonitor) trace(url, method string) *httptracer.TracerResult {
	client := e.tracerPool.Get().(httptracer.HttpTracer)
	report := client.Trace(url, method)
	e.tracerPool.Put(client)
	return report
}

func (um *URLMonitor) getReport(url string) (*httptracer.TracerResult, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	site, ok := um.sites[url]
	if ok {
		site.Update()
		if site.report == nil {
			report := um.trace(url, "GET")
			return report, true
		}
		return site.report, true

	} else {
		return nil, false
	}
}

func (um *URLMonitor) addSite(url string) {
	um.mu.Lock()
	if _, ok := um.sites[url]; ok {
		return
	}

	s := newSite(url, um)

	s.Interval = um.interval
	um.sites[url] = s
	um.mu.Unlock()

	go s.Monitoring()
}

func (um *URLMonitor) clear() {
	log.Debugf("Start clearning")
	var count int
	um.mu.Lock()
	defer um.mu.Unlock()
	for url, site := range um.sites {
		if time.Since(site.GetUpdated()) >= um.clear_interval {
			count++
			site.Stop()
			delete(um.sites, url)
		}
	}
	log.Infof("Clearning %d sites", count)
}

func (um *URLMonitor) clearCron() {
	for range um.clearTicker.C {
		um.clear()
	}
}

func (um *URLMonitor) Stop() {
	um.mu.Lock()
	defer um.mu.Unlock()

	um.clearTicker.Stop()

	for _, site := range um.sites {
		site.Stop()
	}
}
