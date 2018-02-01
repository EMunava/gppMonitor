package eodLog

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type testEodFile struct {
}

func makeEodLogTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(testEodFile)
		s.RetrieveEDOLog()
		return nil, nil
	}
}
