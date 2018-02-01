package monitor

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/eodLog"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
	"github.com/jasonlvhit/gocron"
)

type Service interface {
}

type service struct {
	dateroloverService daterollover.Service
	eodLogService      eodLog.Service
	scheduleBatch      waitSchduleBatch.Service
}

func NewService(dateroloverService daterollover.Service, eodLogService eodLog.Service, scheduleBatch waitSchduleBatch.Service) Service {
	s := &service{dateroloverService: dateroloverService, eodLogService: eodLogService, scheduleBatch: scheduleBatch}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {

	sel := gocron.NewScheduler()

	sel.Every(1).Day().At("23:30").Do(s.dateroloverService.ConfirmDateRollOver)
	sel.Every(1).Day().At("00:30").Do(s.dateroloverService.ConfirmDateRollOver)
	sel.Every(1).Day().At("01:10").Do(s.eodLogService.RetrieveEDOLog)
	sel.Every(1).Day().At("01:30").Do(s.dateroloverService.ConfirmDateRollOver)
	sel.Every(1).Day().At("00:35").Do(s.scheduleBatch.ConfirmWaitSchedSubBatch)

	gocron.NextRun()

	<-sel.Start()
}
