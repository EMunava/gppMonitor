package eodLog

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger, nil)

	eodfile := kithttp.NewServer(makeEodLogTestEndpoint(service), decodeTestEodFile, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/eodfile", eodfile).Methods("GET")

	return r

}

func decodeTestEodFile(_ context.Context, r *http.Request) (interface{}, error) {
	return testEodFile{}, nil
}
