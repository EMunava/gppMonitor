package extractFooterTransactions

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

	transactionSAP := kithttp.NewServer(makeSAPTransactionsEndpoint(service), decodeTestSAPTransaction, gokit.EncodeResponse, opts...)
	transactionLEG := kithttp.NewServer(makeLEGTransactionsEndpoint(service), decodeTestLEGTransaction, gokit.EncodeResponse, opts...)
	transactionLEGSAP := kithttp.NewServer(makeLEGSAPTransactionsEndpoint(service), decodeTestLEGSAPTransaction, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/SAPTransactions", transactionSAP).Methods("GET")
	r.Handle("/LEGTransactions", transactionLEG).Methods("GET")
	r.Handle("/LEGSAPTransactions", transactionLEGSAP).Methods("GET")

	return r

}

func decodeTestSAPTransaction(_ context.Context, r *http.Request) (interface{}, error) {
	return transactionSAP{}, nil
}

func decodeTestLEGTransaction(_ context.Context, r *http.Request) (interface{}, error) {
	return transactionLEG{}, nil
}

func decodeTestLEGSAPTransaction(_ context.Context, r *http.Request) (interface{}, error) {
	return transactionLEGSAP{}, nil
}
