package provider

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type MitakeProvider struct {
	baseURL string
	account SMSAccount
}

func (p *MitakeProvider) SendSMS(phone, message string) {
}

func (p *MitakeProvider) Credit() (int, error) {
	client := resty.New().
		SetHostURL(p.baseURL)

	resp, err := client.R().
		SetFormData(p.Account()).
		Post("/SmQuery")

	if (err != nil) || (resp.StatusCode() != http.StatusOK) {
		return 0, errors.New("查詢餘額失敗")
	}

	scanner := bufio.NewScanner(bytes.NewReader(resp.Body()))
	scanner.Scan()
	line := scanner.Text()

	if line == "statuscode=e" {
		scanner.Scan()
		errorMsg := scanner.Text()
		return 0, errors.New(strings.Split(errorMsg, "=")[1])
	}

	successMsg := strings.Split(line, "=")
	credit, _ := strconv.Atoi(successMsg[1])

	return credit, nil
}

func (p *MitakeProvider) Account() map[string]string {
	return map[string]string{
		"username": p.account.Username,
		"password": p.account.Password,
	}
}
