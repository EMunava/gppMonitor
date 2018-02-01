package gppSelenium

import (
	"github.com/go-kit/kit/metrics"
	"github.com/tebeka/selenium"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) WaitForWaitFor() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "WaitForWaitFor").Add(1)
		s.requestLatency.With("method", "WaitForWaitFor").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.WaitForWaitFor()
}
func (s *instrumentingService) LogIn() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "LogIn").Add(1)
		s.requestLatency.With("method", "LogIn").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.LogIn()
}
func (s *instrumentingService) LogOut() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "LogOut").Add(1)
		s.requestLatency.With("method", "LogOut").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.LogOut()
}

// Override
func (s *instrumentingService) HandleSeleniumError(internal bool, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "HandleSeleniumError").Add(1)
		s.requestLatency.With("method", "HandleSeleniumError").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.HandleSeleniumError(internal, err)
}
func (s *instrumentingService) Driver() selenium.WebDriver {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Driver").Add(1)
		s.requestLatency.With("method", "Driver").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Driver()
}
func (s *instrumentingService) ClickByClassName(cn string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByClassName").Add(1)
		s.requestLatency.With("method", "ClickByClassName").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.ClickByClassName(cn)

}
func (s *instrumentingService) ClickByXPath(xp string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByXPath").Add(1)
		s.requestLatency.With("method", "ClickByXPath").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.ClickByXPath(xp)

}
func (s *instrumentingService) ClickByCSSSelector(cs string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByCSSSelector").Add(1)
		s.requestLatency.With("method", "ClickByCSSSelector").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.ClickByCSSSelector(cs)

}

func (s *instrumentingService) WaitFor(findBy, selector string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "WaitFor").Add(1)
		s.requestLatency.With("method", "WaitFor").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.WaitFor(findBy, selector)
}
