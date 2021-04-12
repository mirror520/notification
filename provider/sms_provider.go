package provider

import (
	"errors"
	"os"

	"github.com/jinzhu/configor"
	"github.com/mirror520/sms/model"

	log "github.com/sirupsen/logrus"
)

type ISMSProvider interface {
	Init()
	Profile() *model.SMSProviderProfile
	SendSMS(*model.SMS) (*model.SMSResult, error)
	Credit() (int, error)
}

var SMSProviderPool map[string]ISMSProvider

func Init() {
	logger := log.WithFields(log.Fields{
		"provider": "ISMSProvider",
		"method":   "Init",
	})

	os.Setenv("CONFIGOR_ENV_PREFIX", "SMS")
	configor.Load(&model.Config, "config.yaml")
	config := model.Config

	SMSProviderPool = make(map[string]ISMSProvider)
	for _, profile := range config.Providers {
		p, err := SMSProviderCreateFactory(profile)
		if err != nil {
			logger.Errorln(err.Error())
		}

		p.Init()
		SMSProviderPool[profile.ID] = p

	}

	logger.Infoln("簡訊提供者初始化完成")
}

func SMSProviderCreateFactory(profile model.SMSProviderProfile) (ISMSProvider, error) {
	switch profile.Type {
	case model.Every8D:
		return &Every8DProvider{
			baseURL: model.Config.Every8D.BaseURL,
			profile: &profile,
		}, nil

	case model.Mitake:
		return &MitakeProvider{
			baseURL: model.Config.Mitake.BaseURL,
			profile: &profile,
		}, nil
	}

	return nil, errors.New("無法提供此簡訊供應商")
}
