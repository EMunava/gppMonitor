package monitor

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/eodLog"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/postingException"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
	"github.com/jasonlvhit/gocron"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/extractFooterTransactions"
)

type Service interface {
}

type service struct {
	dateroloverService daterollover.Service
	eodLogService      eodLog.Service
	scheduleBatch      waitSchduleBatch.Service
	postex             postingException.Service
	transactionService extractFooterTransactions.Service
}

//NewService function creates instances of required external service structs for local use
func NewService(dateroloverService daterollover.Service, eodLogService eodLog.Service, scheduleBatch waitSchduleBatch.Service, postex postingException.Service, transactionService extractFooterTransactions.Service) Service {
	s := &service{dateroloverService: dateroloverService, eodLogService: eodLogService, scheduleBatch: scheduleBatch, postex: postex, transactionService: transactionService}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {

	sel := gocron.NewScheduler()

	sel.Every(1).Day().At("08:00").Do(s.postex.ConfirmPostingException)
	sel.Every(1).Day().At("19:00").Do(s.scheduleBatch.ConfirmWaitSchedSubBatch)
	sel.Every(1).Day().At("23:32").Do(s.dateroloverService.ConfirmDateRollOver)
	sel.Every(1).Day().At("00:22").Do(s.dateroloverService.ConfirmDateRollOver)
	sel.Every(1).Day().At("00:20").Do(s.transactionService.RetrieveLEGSAPTransactions)
	sel.Every(1).Day().At("01:10").Do(s.eodLogService.RetrieveEDOLog)

	gocron.NextRun()

	<-sel.Start()
}
