package extractFooterTransactions

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

func (s *instrumentingService) RetrieveSAPTransactions() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RetrieveSAPTransactions").Add(1)
		s.requestLatency.With("method", "RetrieveSAPTransactions").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RetrieveSAPTransactions()
}

func (s *instrumentingService) RetrieveLEGTransactions() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RetrieveLEGTransactions").Add(1)
		s.requestLatency.With("method", "RetrieveLEGTransactions").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RetrieveLEGTransactions()
}

func (s *instrumentingService) RetrieveLEGSAPTransactions() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RetrieveLEGSAPTransactions").Add(1)
		s.requestLatency.With("method", "RetrieveLEGSAPTransactions").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RetrieveLEGSAPTransactions()
}
