package provider

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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

type mitakeProvider struct {
	log     *zap.Logger
	profile conf.Provider
}

func NewMitakeProvider(profile conf.Provider) message.Provider {
	return &mitakeProvider{
		log:     zap.L().With(zap.String("provider", "mitake")),
		profile: profile,
	}
}

func (p *mitakeProvider) Send(msg *message.Message) error {
	log := p.log.With(
		zap.String("action", "send"),
		zap.String("recipient", msg.Recipient),
	)

	client := resty.New().
		SetBaseURL(p.profile.BaseURL)

	data := map[string]string{
		"username": p.profile.Username,
		"password": p.profile.Password,
		"dstaddr":  msg.Recipient,
		"smbody":   msg.Content,
		"destname": msg.Comment,
		"response": p.profile.CallbackURL,
	}

	resp, err := client.R().
		SetQueryString("CharsetURL=UTF-8").
		SetFormData(data).
		Post("/SmSend")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("send failed")
	}

	var (
		msgID  string
		credit float64
	)

	scanner := bufio.NewScanner(bytes.NewReader(resp.Body()))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "statuscode"):
			statusCode := strings.TrimPrefix(line, "statuscode=")

			if !strings.ContainsAny(statusCode, "0124") {
				scanner.Scan()
				errorMsg := strings.TrimPrefix(scanner.Text(), "Error=")
				return errors.New(errorMsg)
			}

		case strings.HasPrefix(line, "msgid"):
			msgID = strings.TrimPrefix(line, "msgid=")

		case strings.HasPrefix(line, "AccountPoint"):
			accountPointStr := strings.TrimPrefix(line, "AccountPoint=")
			accountPoint, err := strconv.ParseFloat(accountPointStr, 64)
			if err != nil {
				return err
			}

			credit = accountPoint
		}
	}

	log.Info("sent successfully",
		zap.String("msgid", msgID),
		zap.Float64("credit", credit),
	)

	return nil
}

func (p *mitakeProvider) Credit() (float64, error) {
	log := p.log.With(zap.String("action", "credit"))

	client := resty.New().
		SetBaseURL(p.profile.BaseURL)

	data := map[string]string{
		"username": p.profile.Username,
		"password": p.profile.Password,
	}

	resp, err := client.R().
		SetFormData(data).
		Post("/SmQuery")

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, errors.New("query failed")
	}

	scanner := bufio.NewScanner(bytes.NewReader(resp.Body()))
	scanner.Scan()
	line := scanner.Text()

	if strings.TrimPrefix(line, "statuscode=") == "e" {
		scanner.Scan()
		errorMsg := strings.TrimPrefix(scanner.Text(), "Error=")
		return 0, errors.New(errorMsg)
	}

	accountPoint := strings.TrimPrefix(line, "AccountPoint=")
	credit, err := strconv.ParseFloat(accountPoint, 64)
	if err != nil {
		return 0, err
	}

	log.Info("credit updated", zap.Float64("credit", credit))
	return credit, nil
}

func (p *mitakeProvider) Callback(query url.Values) (string, error) {
	log := p.log.With(
		zap.String("action", "callback"),
	)

	mid := query.Get("msgid")
	if mid == "" {
		return "", errors.New("invalid id")
	}

	recipient := query.Get("dstaddr")
	recipient = strings.Replace(recipient, "09", "+8869", 1)

	status := query.Get("statusstr")

	sendTime, err := time.ParseInLocation("20060102150405", query.Get("dlvtime"), time.Local)
	if err != nil {
		return "", err
	}

	receiveTime, err := time.ParseInLocation("20060102150405", query.Get("donetime"), time.Local)
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

	response := fmt.Sprintf("magicid=sms_gateway_rpack\nmsgid=%s\n", mid)
	return response, nil
}
