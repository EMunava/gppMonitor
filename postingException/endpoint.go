package postingException

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type postingExceptionTest struct {
}

func makePostingExceptionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(postingExceptionTest)
		s.ConfirmPostingException()
		return nil, nil
	}
}
