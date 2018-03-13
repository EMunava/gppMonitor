package transactionCountGUI

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

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
	s.Service.ConfirmWaitSchedSubBatch()
}
