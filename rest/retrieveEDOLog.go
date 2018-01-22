package rest

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"net/http"
)

func retrieveEDOLog(w http.ResponseWriter, r *http.Request) {

	sftp.RetrieveEDOLog()
	w.WriteHeader(http.StatusOK)
}

