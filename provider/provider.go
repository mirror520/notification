package provider

import (
	"errors"

	"github.com/mirror520/notification/conf"
	"github.com/mirror520/notification/message"
)

func NewProvider(profile conf.Provider) (message.Provider, error) {
	switch profile.Type {
	case conf.Every8D:
		return NewEvery8DProvider(profile), nil
	case conf.Mitake:
		return NewMitakeProvider(profile), nil
	default:
		return nil, errors.New("provider not found")
	}
}
