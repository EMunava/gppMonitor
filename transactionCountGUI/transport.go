package transactionCountGUI

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

	waitSchedSubBatch := kithttp.NewServer(makeWaitBatchTestEndpoint(service), decodeTestBatch, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/waitschedulebatch", waitSchedSubBatch).Methods("GET")

	return r

}

func decodeTestBatch(_ context.Context, r *http.Request) (interface{}, error) {
	return waitBatchTest{}, nil
}
