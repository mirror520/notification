package provider

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mirror520/sms/model"

	log "github.com/sirupsen/logrus"
)

type MitakeProvider struct {
	baseURL string
	account model.SMSAccount
	credit  int
}

func (p *MitakeProvider) Init() {
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

func (p *MitakeProvider) SendSMS(sms *model.SMS) (*model.SMSResult, error) {
	client := resty.New().
		SetHostURL(p.baseURL)

	resp, err := client.R().
		SetQueryString("CharsetURL=UTF-8").
		SetFormData(p.AccountWithSMS(sms)).
		Post("/SmSend")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return nil, errors.New("傳送簡訊失敗")
	}

	var result = &model.SMSResult{}
	scanner := bufio.NewScanner(bytes.NewReader(resp.Body()))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "statuscode"):
			statusCode := strings.TrimPrefix(line, "statuscode=")

			if !strings.ContainsAny(statusCode, "0124") {
				scanner.Scan()
				errorMsg := strings.TrimPrefix(scanner.Text(), "Error=")
				return nil, errors.New(errorMsg)
			}

		case strings.HasPrefix(line, "msgid"):
			msgID := strings.TrimPrefix(line, "msgid=")
			result.ID = msgID

		case strings.HasPrefix(line, "AccountPoint"):
			accountPoint := strings.TrimPrefix(line, "AccountPoint=")
			credit, _ := strconv.Atoi(accountPoint)
			result.Credit = credit
		}
	}

	return result, nil
}

func (p *MitakeProvider) Credit() (int, error) {
	client := resty.New().
		SetHostURL(p.baseURL)

	resp, err := client.R().
		SetFormData(p.Account()).
		Post("/SmQuery")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return -1, errors.New("查詢餘額失敗")
	}

	scanner := bufio.NewScanner(bytes.NewReader(resp.Body()))
	scanner.Scan()
	line := scanner.Text()

	if strings.TrimPrefix(line, "statuscode=") == "e" {
		scanner.Scan()
		errorMsg := strings.TrimPrefix(scanner.Text(), "Error=")
		return -1, errors.New(errorMsg)
	}

	accountPoint := strings.TrimPrefix(line, "AccountPoint=")
	credit, _ := strconv.Atoi(accountPoint)

	return credit, nil
}

func (p *MitakeProvider) Account() map[string]string {
	return map[string]string{
		"username": p.account.Username,
		"password": p.account.Password,
	}
}

func (p *MitakeProvider) AccountWithSMS(sms *model.SMS) map[string]string {
	account := p.Account()
	account["dstaddr"] = sms.Phone
	account["smbody"] = sms.Message
	account["destname"] = sms.Comment
	return account
}
