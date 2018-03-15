package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/remoteTelegramCommands"
	"github.com/weAutomateEverything/gppMonitor/daterollover"
	"github.com/weAutomateEverything/gppMonitor/eodLog"
	"github.com/weAutomateEverything/gppMonitor/gppSelenium"
	"github.com/weAutomateEverything/gppMonitor/postingException"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"github.com/weAutomateEverything/gppMonitor/transactionCountGUI"
	"github.com/weAutomateEverything/gppMonitor/transactionCountLog"
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

	calloutService := callout.NewKubernetesCalloutProxy("")

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

	dateRolloverService := daterollover.NewService(calloutService, alertService, gppSeleniumService, remoteTelegramService)
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

	eodLogService := eodLog.NewService(calloutService, sftpService, alertService)
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

	waitScheduleBatchService := transactionCountGUI.NewService(alertService, gppSeleniumService)
	waitScheduleBatchService = transactionCountGUI.NewLoggingService(log.With(logger, "component", "wait_schedule_batch"), waitScheduleBatchService)
	waitScheduleBatchService = transactionCountGUI.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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

	postingExceptionService := postingException.NewService(alertService, gppSeleniumService)
	postingExceptionService = postingException.NewLoggingService(log.With(logger, "component", "wait_schedule_batch"), postingExceptionService)
	postingExceptionService = postingException.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "posting_exception",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "posting_exception",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), postingExceptionService)
	transactionService := transactionCountLog.NewService(calloutService, sftpService, alertService)
	transactionService = transactionCountLog.NewLoggingService(log.With(logger, "component", "Transaction Count"), transactionService)
	transactionService = transactionCountLog.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "transaction_count",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "transaction_count",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), transactionService)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/daterollover", daterollover.MakeHandler(dateRolloverService, httpLogger))
	mux.Handle("/eodfile", eodLog.MakeHandler(eodLogService, httpLogger))
	mux.Handle("/waitschedulebatch", transactionCountGUI.MakeHandler(waitScheduleBatchService, httpLogger))
	mux.Handle("/postingException", postingException.MakeHandler(postingExceptionService, httpLogger))
	mux.Handle("/SAPTransactions", transactionCountLog.MakeHandler(transactionService, httpLogger))
	mux.Handle("/LEGTransactions", transactionCountLog.MakeHandler(transactionService, httpLogger))
	mux.Handle("/LEGSAPTransactions", transactionCountLog.MakeHandler(transactionService, httpLogger))

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
