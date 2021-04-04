package provider

import (
	"errors"

	"github.com/mirror520/sms/environment"
)

type SMSProvider interface {
	SendSMS(phone, message string)
	Credit() (int, error)
}

type SMSProviderType int

const (
	Every8D SMSProviderType = iota
	Mitake
)

type SMSAccount struct {
	Username string
	Password string
}

func SMSProviderCreateFactory(provider SMSProviderType) (SMSProvider, error) {
	switch provider {
	case Every8D:
		return &Every8DProvider{
			baseURL: environment.EVERY8D_BASEURL,
			account: SMSAccount{
				Username: environment.EVERY8D_USERNAME,
				Password: environment.EVERY8D_PASSWORD,
			},
		}, nil

	case Mitake:
		return &MitakeProvider{
			baseURL: environment.MITAKE_BASEURL,
			account: SMSAccount{
				Username: environment.MITAKE_USERNAME,
				Password: environment.MITAKE_PASSWORD,
			},
		}, nil
	}

	return nil, errors.New("無法提供此簡訊供應商")
}