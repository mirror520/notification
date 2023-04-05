package notification

import (
	"context"
	"errors"
	"net/url"

	"github.com/go-kit/kit/endpoint"

	"github.com/mirror520/notification/message"
)

func SendEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		msg, ok := request.(*message.Message)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := svc.Send(msg); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func CreditEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Credit(req)
	}
}

type CallbackRequest struct {
	Provider string
	Values   url.Values
}

func CallbackEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(*CallbackRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Callback(req.Values, req.Provider)
	}
}
