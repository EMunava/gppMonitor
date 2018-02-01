package rest

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/selenium"
	"net/http"
)

func runDateRolloverCheck(w http.ResponseWriter, r *http.Request) {

	selenium.CallSeleniumDateCheck()
	w.WriteHeader(http.StatusOK)
}

func runWaitSchedBatchCheck(w http.ResponseWriter, r *http.Request) {

	selenium.CallWaitSchedBatchCheck()
	w.WriteHeader(http.StatusOK)
}
