package extractFooterTransactions

import (
	"time"
	"log"
)

type loggingService struct {
	logger log.Logger
	Service
}

//NewLoggingService function creates an instance of the loggingService struct for local use
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) WaitForWaitFor() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ConfirmWaitSchedSubBatch",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RetrieveLEGSAPTransactions()
	s.Service.RetrieveLEGTransactions()
	s.Service.RetrieveSAPTransactions()
}
