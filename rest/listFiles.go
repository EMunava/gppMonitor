package rest

import (
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"net/http"
)

func listFiles(w http.ResponseWriter, r *http.Request) {

	sftp.ListFiles()
	w.WriteHeader(http.StatusOK)
}
