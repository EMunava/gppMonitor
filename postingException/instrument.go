package postingException

import (
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

//NewInstrumentService function creates an instance of the instrumentingService struct for local use
func NewInstrumentService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) ConfirmPostingException() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ConfirmPostingException").Add(1)
		s.requestLatency.With("method", "ConfirmPostingException").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ConfirmPostingException()
}
