package provider

import (
	"errors"

	"github.com/mirror520/sms/model"
)

type ISMSProvider interface {
	SendSMS(*model.SMS) (*model.SMSResult, error)
	Credit() (int, error)
}

func SMSProviderCreateFactory(p model.SMSProvider) (ISMSProvider, error) {
	switch p.Type {
	case model.Every8D:
		return &Every8DProvider{
			baseURL: model.Config.Every8D.BaseURL,
			account: p.Account,
		}, nil

	case model.Mitake:
		return &MitakeProvider{
			baseURL: model.Config.Mitake.BaseURL,
			account: p.Account,
		}, nil
	}

	return nil, errors.New("無法提供此簡訊供應商")
}
