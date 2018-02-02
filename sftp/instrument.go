package sftp

import (
	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) GetFilesInPath(path string) ([]File, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetFilesInPath").Add(1)
		s.requestLatency.With("method", "GetFilesInPath").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetFilesInPath(path)
}

func (s *instrumentingService) RetrieveFile(path, file string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RetrieveFile").Add(1)
		s.requestLatency.With("method", "RetrieveFile").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RetrieveFile(path, file)
}
