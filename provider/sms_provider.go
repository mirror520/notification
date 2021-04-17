package provider

import (
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/mirror520/sms/model"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

type ISMSProvider interface {
	Init()
	Profile() *model.SMSProviderProfile
	SendSMS(*model.SMS) (*model.SMSResult, error)
	Credit() (int, error)
	Callback(*url.Values) (string, string, error)
}

var InfluxDB influxdb.Client

var (
	providerPool   map[string]ISMSProvider
	masterProvider ISMSProvider

	timeLayout   string
	timeLocation *time.Location
)

func Init() {
	logger := log.WithFields(log.Fields{
		"provider": "ISMSProvider",
		"method":   "Init",
	})

	timeLayout = "20060102150405"
	timeLocation, _ = time.LoadLocation("Asia/Taipei")

	os.Setenv("CONFIGOR_ENV_PREFIX", "SMS")
	configor.Load(&model.Config, "config.yaml")
	config := model.Config

	providerPool = make(map[string]ISMSProvider)
	for _, profile := range config.Providers {
		p, err := SMSProviderCreateFactory(profile)
		if err != nil {
			logger.Errorln(err.Error())
		}

		if profile.Role == model.Master {
			masterProvider = p
		}

		p.Init()
		providerPool[profile.ID] = p
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
			baseURL:     model.Config.Mitake.BaseURL,
			callbackURL: model.Config.Mitake.CallbackURL,
			profile:     &profile,
		}, nil
	}

	return nil, errors.New("無法產生此簡訊提供商")
}

func SwitchSMSProviderToMaster(master ISMSProvider) {
	masterProvider = master

	for _, p := range providerPool {
		if p.Profile().ID == master.Profile().ID {
			p.Profile().Role = model.Master
		} else {
			p.Profile().Role = model.Backup
		}
	}
}

func NewSMSStatusToTSDB(pid, status string, delay float64, sendTime time.Time) {
	config := model.Config
	writeAPI := InfluxDB.WriteAPI(config.InfluxDB.Org, config.InfluxDB.Bucket)

	p := influxdb.NewPointWithMeasurement("send_status").
		AddTag("pid", pid).
		AddTag("status", status).
		AddField("delay", delay).
		SetTime(sendTime)

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

func SMSProvider(pid string) (ISMSProvider, error) {
	p, ok := providerPool[pid]

	if !ok {
		return nil, errors.New("無此簡訊提供商")
	}

	return p, nil
}

func SMSMasterProvider() ISMSProvider {
	return masterProvider
}
