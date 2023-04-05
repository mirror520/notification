package provider

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/mirror520/notification/conf"
	"github.com/mirror520/notification/message"
)

type every8DProvider struct {
	log     *zap.Logger
	profile conf.Provider
}

func NewEvery8DProvider(profile conf.Provider) message.Provider {
	return &every8DProvider{
		log:     zap.L().With(zap.String("provider", "every8d")),
		profile: profile,
	}
}

func (p *every8DProvider) Send(msg *message.Message) error {
	log := p.log.With(
		zap.String("action", "send"),
		zap.String("recipient", msg.Recipient),
	)

	client := resty.New().
		SetBaseURL(p.profile.BaseURL)

	data := map[string]string{
		"UID":  p.profile.Username,
		"PWD":  p.profile.Password,
		"DEST": msg.Recipient,
		"MSG":  msg.Content,
		"SB":   msg.Comment,
	}

	resp, err := client.R().
		SetFormData(data).
		Post("/sendSMS.ashx")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("send failed")
	}

	body := string(resp.Body())
	contents := strings.Split(body, ",")

	credit, err := strconv.ParseFloat(contents[0], 64)
	if err != nil {
		return err
	}

	if credit < 0 {
		return errors.New(strings.Trim(contents[1], " "))
	}

	batchID := contents[4]

	log.Info("sent successfully",
		zap.String("msgid", batchID),
		zap.Float64("credit", credit),
	)
	return nil
}

func (p *every8DProvider) Credit() (float64, error) {
	log := p.log.With(
		zap.String("action", "credit"),
	)

	client := resty.New().
		SetBaseURL(p.profile.BaseURL)

	data := map[string]string{
		"UID": p.profile.Username,
		"PWD": p.profile.Password,
	}

	resp, err := client.R().
		SetFormData(data).
		Post("/getCredit.ashx")

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, errors.New("query failed")
	}

	body := string(resp.Body())
	contents := strings.Split(body, ",")

	credit, err := strconv.ParseFloat(contents[0], 64)
	if err != nil {
		return 0, err
	}

	if credit < 0 {
		return 0, errors.New(strings.TrimSpace(contents[1]))
	}

	log.Info("credit updated", zap.Float64("credit", credit))
	return credit, nil
}

func (p *every8DProvider) Callback(query url.Values) (string, error) {
	log := p.log.With(
		zap.String("action", "callback"),
	)

	mid := query.Get("BatchID")
	if mid == "" {
		return "", errors.New("invalid id")
	}

	recipient := query.Get("RM")

	status := query.Get("STATUS")

	sendTime, err := time.ParseInLocation("20060102150405", query.Get("ST"), time.Local)
	if err != nil {
		return "", err
	}

	receiveTime, err := time.ParseInLocation("20060102150405", query.Get("RT"), time.Local)
	if err != nil {
		return "", err
	}

	delay := receiveTime.Sub(sendTime)

	// TODO: metrics
	log.Info("metrics",
		zap.String("msgid", mid),
		zap.String("recipient", recipient),
		zap.String("status", status),
		zap.Time("send_time", sendTime),
		zap.Time("receive_time", receiveTime),
		zap.Duration("delay", delay),
	)

	return "ok", nil
}
