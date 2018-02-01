package daterollover

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)

	daterollover := kithttp.NewServer(makeTestDateRolloverEndpoint(service), decodeTestDateRollover, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/daterollover", daterollover).Methods("GET")

	return r

}

func decodeTestDateRollover(_ context.Context, r *http.Request) (interface{}, error) {
	return testDateRollOver{}, nil
}
