package provider

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mirror520/sms/model"

	log "github.com/sirupsen/logrus"
)

type Every8DProvider struct {
	baseURL string
	profile *model.SMSProviderProfile
	credit  int
}

func (p *Every8DProvider) Init() {
	logger := log.WithFields(log.Fields{
		"provider": "Every8DProvider",
		"method":   "Init",
	})

	credit, err := p.Credit()
	if err != nil {
		logger.Errorln(err.Error())
	}

	p.credit = credit
	logger.Infoln("初始化完成")
}

func (p *Every8DProvider) SendSMS(sms *model.SMS) (*model.SMSResult, error) {
	client := resty.New().
		SetHostURL(p.baseURL)

	resp, err := client.R().
		SetFormData(p.AccountWithSMS(sms)).
		Post("/sendSMS.ashx")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return nil, errors.New("傳送簡訊失敗")
	}

	body := string(resp.Body())
	contents := strings.Split(body, ",")

	credit, _ := strconv.ParseFloat(contents[0], 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	if credit < 0 {
		return nil, errors.New(strings.Trim(contents[1], " "))
	}

	batchID := contents[4]
	result := &model.SMSResult{
		ID:     batchID,
		Credit: int(credit),
	}

	return result, nil
}

func (p *Every8DProvider) Credit() (int, error) {
	client := resty.New().
		SetHostURL(p.baseURL)

	resp, err := client.R().
		SetFormData(p.Account()).
		Post("/getCredit.ashx")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return -1, errors.New("查詢餘額失敗")
	}

	body := string(resp.Body())
	contents := strings.Split(body, ",")

	credit, _ := strconv.Atoi(contents[0])
	if credit < 0 {
		return -1, errors.New(strings.Trim(contents[1], " "))
	}

	return credit, nil
}

func (p *Every8DProvider) Profile() *model.SMSProviderProfile {
	return p.profile
}

func (p *Every8DProvider) Account() map[string]string {
	return map[string]string{
		"UID": p.Profile().Account.Username,
		"PWD": p.Profile().Account.Password,
	}
}

func (p *Every8DProvider) AccountWithSMS(sms *model.SMS) map[string]string {
	account := p.Account()
	account["DEST"] = sms.Phone
	account["MSG"] = sms.Message
	account["SB"] = sms.Comment
	return account
}
