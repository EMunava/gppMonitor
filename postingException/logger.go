package postingException

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

//NewLoggingService function creates an instance of the loggingService struct for local use
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) ConfirmWaitPostingExceptions() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ConfirmWaitPostingExceptions",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ConfirmPostingException()
}
