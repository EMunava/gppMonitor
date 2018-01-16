package rest

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/selenium"
	"net/http"
)

func runSelenium(w http.ResponseWriter, r *http.Request) {

	selenium.CallSeleniumDateCheck()
	w.WriteHeader(http.StatusOK)
}
