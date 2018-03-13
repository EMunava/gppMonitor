package transactionCountGUI

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type waitBatchTest struct {
}

func makeWaitBatchTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(waitBatchTest)
		s.ExtractTransactionCount()
		return nil, nil
	}
}
