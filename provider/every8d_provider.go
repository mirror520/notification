package provider

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

var baseURL = "https://oms.every8d.com/API21/HTTP"

type Every8DProvider struct {
	account SMSAccount
}

func (p *Every8DProvider) SendSMS(phone, message string) {

}

func (p *Every8DProvider) Credit() (int, error) {
	client := resty.New().
		SetHostURL(baseURL)

	resp, err := client.R().
		SetFormData(p.Account()).
		Post("/getCredit.ashx")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return 0, errors.New("查詢餘額失敗")
	}

	contents := strings.Split(string(resp.Body()), ",")

	credit, _ := strconv.Atoi(contents[0])
	if credit < 0 {
		return 0, errors.New(strings.Trim(contents[1], " "))
	}

	return credit, nil
}

func (p *Every8DProvider) Account() map[string]string {
	return map[string]string{
		"UID": p.account.Username,
		"PWD": p.account.Password,
	}
}
