package main

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/monitor"
	"github.com/zamedic/go2hal/alert"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/go-kit/kit/log"
)

func main() {

	var logger log.Logger

	alertService := alert.NewKubernetesAlertProxy("")
	_ = monitor.NewService(alertService)

	errs := make(chan error, 2)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)


}
