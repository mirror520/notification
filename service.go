package notification

import (
	"errors"
	"net/url"

	"go.uber.org/zap"

	"github.com/mirror520/notification/message"
)

type Service interface {
	Send(msg *message.Message) error
	Credit(provider string) (float64, error)
	Callback(values url.Values, provider string) (string, error)
}

type ServiceMiddleware func(Service) Service

func NewService(providers map[string]message.Provider) Service {
	return &service{
		zap.L().With(zap.String("service", "notification")),
		providers,
	}
}

type service struct {
	log       *zap.Logger
	providers map[string]message.Provider
}

func (svc *service) Send(msg *message.Message) error {
	log := svc.log.With(
		zap.String("action", "send"),
		zap.String("recipient", msg.Recipient),
	)

	count := 0
	for name, p := range svc.providers {
		if err := p.Send(msg); err != nil {
			log.Error(err.Error(), zap.String("provider", name))
			continue
		}

		count++
	}

	if count == 0 {
		return errors.New("send failed")
	}

	return nil
}

func (svc *service) Credit(provider string) (float64, error) {
	p, ok := svc.providers[provider]
	if !ok {
		return 0, errors.New("provider not found")
	}

	return p.Credit()
}

func (svc *service) Callback(values url.Values, provider string) (string, error) {
	p, ok := svc.providers[provider]
	if !ok {
		return "", errors.New("provider not found")
	}

	return p.Callback(values)
}
