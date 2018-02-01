package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/zamedic/go2hal/alert"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/eodLog"
)

type Service interface {
}

type service struct {
	alert alert.Service
}

func NewService(alert alert.Service) Service {
	s := &service{alert: alert}
	s.schedule()
	return s
}

func (s *service) schedule() {
	sftpService := sftp.NewService()

	sel := gocron.NewScheduler()
	dateRollover := func() { daterollover.NewService(s.alert).ConfirmDateRollOver() }
	scheduleBatch := func() { waitSchduleBatch.NewService(s.alert).ConfirmWaitSchedSubBatch() }
	eofFileCheck := func() { eodLog.NewService(sftpService, s.alert).RetrieveEDOLog() }

	sel.Every(1).Day().At("23:30").Do(dateRollover)
	sel.Every(1).Day().At("00:30").Do(dateRollover)
	sel.Every(1).Day().At("01:10").Do(eofFileCheck)
	sel.Every(1).Day().At("01:30").Do(dateRollover)
	sel.Every(1).Day().At("00:35").Do(scheduleBatch)

	gocron.NextRun()

	<-sel.Start()
}
