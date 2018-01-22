package rest

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"net/http"
)

func retreiveEDOLog(w http.ResponseWriter, r *http.Request) {

	sftp.RetreiveEDOLog()
	w.WriteHeader(http.StatusOK)
}

