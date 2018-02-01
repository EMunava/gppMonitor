package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
)

type Service interface{

}

type service struct{
	dateRollover daterollover.Service

}


func NewService() Service {

}

func (s *service)schedule(){
	sel := gocron.NewScheduler()
	sel.Every(1).Day().At("23:30").Do(seleniumDateRolloverCheck)
	sel.Every(1).Day().At("00:30").Do(seleniumDateRolloverCheck)
	sel.Every(1).Day().At("01:30").Do(seleniumDateRolloverCheck)
	gocron.NextRun()

	<-sel.Start()
}