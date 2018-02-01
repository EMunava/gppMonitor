package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/zamedic/go2hal/alert"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
)

type Service interface {
}

type service struct {
	alert alert.Service
}

func NewService(alert alert.Service)Service{
	return &service{alert:alert}
}


func (s *service) schedule() {
	sel := gocron.NewScheduler()
	dateRollover := func() { daterollover.NewService(s.alert).ConfirmDateRollOver() }
	scheduleBatch := func() { waitSchduleBatch.NewService(s.alert).ConfirmWaitSchedSubBatch() }
	sel.Every(1).Day().At("23:30").Do(dateRollover)
	sel.Every(1).Day().At("00:30").Do(dateRollover)
	sel.Every(1).Day().At("01:30").Do(dateRollover)

	sel.Every(1).Day().At("00:35").Do(scheduleBatch)

	gocron.NextRun()

	<-sel.Start()
}
