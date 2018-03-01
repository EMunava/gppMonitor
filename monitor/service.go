package monitor

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/eodLog"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/extractFooterTransactions"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/postingException"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
	"github.com/jasonlvhit/gocron"
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

	postex := gocron.NewScheduler()
	postexConst := gocron.NewScheduler()
	confirmWaitSched := gocron.NewScheduler()
	confirmDateRoll := gocron.NewScheduler()
	retreiveSAP := gocron.NewScheduler()
	retreiveLEG := gocron.NewScheduler()
	retreiveLEGSAP := gocron.NewScheduler()
	retreiveEDOLog := gocron.NewScheduler()

	postex.Every(1).Day().At("08:00").Do(s.postex.ConfirmPostingException)
	postexConst.Every(1).Hour().Do(s.postex.ConfirmPostingException)
	confirmWaitSched.Every(1).Day().At("19:00").Do(s.scheduleBatch.ConfirmWaitSchedSubBatch)
	confirmDateRoll.Every(1).Day().At("23:32").Do(s.dateroloverService.ConfirmDateRollOver)
	retreiveSAP.Every(1).Day().At("00:05").Do(s.transactionService.RetrieveSAPTransactions)
	retreiveLEG.Every(1).Day().At("01:32").Do(s.transactionService.RetrieveLEGTransactions)
	retreiveLEGSAP.Every(1).Day().At("00:20").Do(s.transactionService.RetrieveLEGSAPTransactions)
	retreiveEDOLog.Every(1).Day().At("01:10").Do(s.eodLogService.RetrieveEDOLog)

	gocron.NextRun()

	<-postex.Start()
	<-postexConst.Start()
	<-confirmWaitSched.Start()
	<-confirmDateRoll.Start()
	<-retreiveSAP.Start()
	<-retreiveLEG.Start()
	<-retreiveLEGSAP.Start()
	<-retreiveEDOLog.Start()
}
