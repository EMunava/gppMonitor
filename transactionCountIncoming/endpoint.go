package transactionCountIncoming

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type transactionSAP struct{}
type transactionLEG struct{}
type transactionLEGSAP struct{}

func makeSAPTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(transactionSAP)
		s.RetrieveSAPTransactions()
		return nil, nil
	}
}

func makeLEGTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(transactionLEG)
		s.RetrieveLEGTransactions()
		return nil, nil
	}
}

func makeLEGSAPTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(transactionLEGSAP)
		s.RetrieveLEGSAPTransactions()
		return nil, nil
	}
}
