package main

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/daterollover"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/eodLog"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/monitor"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/waitSchduleBatch"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/remoteTelegramCommands"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	fieldKeys := []string{"method"}

	alertService := alert.NewKubernetesAlertProxy("")
	alertService.SendAlert("test")

	remoteTelegramService := remoteTelegramCommands.NewRemoteCommandClientService()

	gppSeleniumService := gppSelenium.NewService(alertService)
	gppSeleniumService = gppSelenium.NewLoggingService(log.With(logger, "component", "selenium"), gppSeleniumService)
	gppSeleniumService = gppSelenium.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "selenium",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "selenium",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), gppSeleniumService)

	dateRolloverService := daterollover.NewService(alertService, gppSeleniumService, remoteTelegramService)
	dateRolloverService = daterollover.NewLoggingService(log.With(logger, "component", "date_rollover"), dateRolloverService)
	dateRolloverService = daterollover.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "date_rollover",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "date_rollover",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), dateRolloverService)

	sftpService := sftp.NewService()
	sftpService = sftp.NewLoggingService(log.With(logger, "component", "sftp"), sftpService)
	sftpService = sftp.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "sftp",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "sftp",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), sftpService)

	eodLogService := eodLog.NewService(sftpService, alertService)
	eodLogService = eodLog.NewLoggingService(log.With(logger, "component", "eod_log"), eodLogService)
	eodLogService = eodLog.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "eod_log",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "eod_log",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), eodLogService)

	waitScheduleBatchService := waitSchduleBatch.NewService(alertService, gppSeleniumService)
	waitScheduleBatchService = waitSchduleBatch.NewLoggingService(log.With(logger, "component", "wait_schedule_batch"), waitScheduleBatchService)
	waitScheduleBatchService = waitSchduleBatch.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "wait_schedule_batch",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "wait_schedule_batch",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), waitScheduleBatchService)

	_ = monitor.NewService(dateRolloverService, eodLogService, waitScheduleBatchService)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/daterollover", daterollover.MakeHandler(dateRolloverService, httpLogger))
	mux.Handle("/eodfile", eodLog.MakeHandler(eodLogService, httpLogger))
	mux.Handle("/waitchedulebatch", waitSchduleBatch.MakeHandler(waitScheduleBatchService, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":8001", "msg", "listening")
		errs <- http.ListenAndServe(":8001", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
