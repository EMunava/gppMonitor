package daterollover

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type testDateRollOver struct {
}

func makeTestDateRolloverEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(testDateRollOver)
		s.ConfirmDateRollOver()
		return nil, nil
	}
}
