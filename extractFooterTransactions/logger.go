package extractFooterTransactions

import (
	"time"
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

//NewLoggingService function creates an instance of the loggingService struct for local use
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) RetreiveSAPTransactions() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RetreiveSapTransactions",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RetrieveSAPTransactions()
}

func (s *loggingService) RetreiveLEGTransactions() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RetreiveLEGTransactions",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RetrieveLEGTransactions()
}

func (s *loggingService) RetreiveLEGSAPTransactions() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RetreiveLEGSAPTransactions",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RetrieveLEGSAPTransactions()
}