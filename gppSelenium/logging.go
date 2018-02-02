package gppSelenium

import (
	"github.com/go-kit/kit/log"
	"github.com/tebeka/selenium"
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
			"method", "WaitForWaitFor",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.WaitForWaitFor()
}
func (s *loggingService) LogIn() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "LogIn",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.LogIn()

}
func (s *loggingService) LogOut() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "LogOut",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.LogOut()
}

// Override
func (s *loggingService) HandleSeleniumError(internal bool, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "HandleSeleniumError",
			"internal", internal,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.HandleSeleniumError(internal, err)

}
func (s *loggingService) Driver() (driver selenium.WebDriver) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Driver",
			"driver", driver,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Driver()

}
func (s *loggingService) ClickByClassName(cn string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ClickByClassName",
			"cn", cn,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ClickByClassName(cn)

}
func (s *loggingService) ClickByXPath(xp string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ClickByXPath",
			"xp", xp,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ClickByXPath(xp)

}
func (s *loggingService) ClickByCSSSelector(cs string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ClickByCSSSelector",
			"cs", cs,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ClickByCSSSelector(cs)

}

func (s *loggingService) WaitFor(findBy, selector string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "WaitFor",
			"findBy", findBy,
			"selector", selector,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.WaitFor(findBy, selector)

}

func (s *loggingService) NewClient() (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "NewClient",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.NewClient()
}
