package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/aredoff/reagate/pkg/httptracer"
)

const (
	DEFAULT_CLEAR_INTERVAL = 120 * time.Second
)

func New() *URLMonitor {
	um := URLMonitor{
		tracerPool: sync.Pool{
			New: func() interface{} { return httptracer.New() },
		},
		sites:       make(map[string]*site),
		mu:          &sync.RWMutex{},
		clearTicker: time.NewTicker(DEFAULT_CLEAR_INTERVAL),
		interval:    DEFAULT_INTERVAL,
	}

	go um.clearCron()

	return &um
}

type URLMonitor struct {
	tracerPool  sync.Pool
	sites       map[string]*site
	interval    time.Duration
	clearTicker *time.Ticker
	mu          *sync.RWMutex
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
	um.addSite(url, DEFAULT_INTERVAL)
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

func (um *URLMonitor) addSite(url string, interval time.Duration) {
	um.mu.Lock()
	if _, ok := um.sites[url]; ok {
		return
	}

	s := newSite(url, um)
	s.Interval = interval
	um.sites[url] = s
	um.mu.Unlock()

	go s.Monitoring()
}

func (um *URLMonitor) clear() {
	um.mu.Lock()
	defer um.mu.Unlock()
	fmt.Println("Start clearning")
	for url, site := range um.sites {
		if time.Since(site.GetUpdated()) >= um.interval*1 {
			site.Stop()
			delete(um.sites, url)
		}
	}
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
