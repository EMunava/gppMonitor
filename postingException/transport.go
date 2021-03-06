package postingException

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"net/http"
)

//MakeHandler creates a new Mux router to handle postingException REST requests
func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger, nil)

	postingException := kithttp.NewServer(makePostingExceptionEndpoint(service), decodeTestBatch, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/postingException", postingException).Methods("GET")

	return r

}

func decodeTestBatch(_ context.Context, r *http.Request) (interface{}, error) {
	return postingExceptionTest{}, nil
}
