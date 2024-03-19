package app

import (
	"math/rand"
	"sync"
	"time"

	"github.com/aredoff/reagate/internal/log"
	"github.com/aredoff/reagate/pkg/httptracer"
)

func newSite(url string, um *URLMonitor) *site {
	return &site{
		um:       um,
		Interval: DEFAULT_INTERVAL,
		URL:      url,
		ticker:   time.NewTicker(um.interval),
		done:     make(chan struct{}),
		mu:       &sync.Mutex{},
		updated:  time.Now(),
	}

}

type site struct {
	um       *URLMonitor
	Interval time.Duration
	URL      string
	ticker   *time.Ticker
	updated  time.Time
	report   *httptracer.TracerResult
	done     chan struct{}
	mu       *sync.Mutex
}

func (s *site) Update() {
	s.mu.Lock()
	s.updated = time.Now()
	s.mu.Unlock()
}

func (s *site) GetUpdated() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updated
}

func (s *site) SetReport(report *httptracer.TracerResult) {
	s.mu.Lock()
	s.report = report
	s.mu.Unlock()
}

func (s *site) GetReport() *httptracer.TracerResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.report
}

func (s *site) Stop() {
	s.ticker.Stop()
	close(s.done)
}

func (s *site) Monitoring() {
	if s.URL == "" {
		log.Fatal("URL is empty")
	}
	if s.Interval == 0 {
		log.Fatal("Interval is 0")
	}

	time.Sleep(time.Duration(rand.Intn(int(s.Interval.Milliseconds())) * int(time.Millisecond)))
	log.Infof("Start monitoring %s", s.URL)
	defer log.Infof("Stop monitorin %s", s.URL)
	for {
		select {
		case <-s.done:
			return
		case <-s.ticker.C:
			s.report = s.um.trace(s.URL, "GET")
		}
	}
}
