package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/sms/model"
	"github.com/stretchr/testify/suite"
)

type SMSTestSuite struct {
	suite.Suite
	router *gin.Engine
	sms    *model.SMS
}

func (suite *SMSTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = setRouter()

	suite.sms = &model.SMS{
		Phone:   os.Getenv("SMS_TESTPHONE"),
		Message: "現在時間: " + time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (suite *SMSTestSuite) TestSMSCreditByEvery8D() {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sms/credit/every8d", nil)
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")
	if res.Code == http.StatusOK {
		var result *model.Result
		json.NewDecoder(res.Body).Decode(&result)

		credit := result.Data.(float64)
		suite.GreaterOrEqual(int(credit), 0, "餘額應大於或等於 0")
	}
}

func (suite *SMSTestSuite) TestSMSCreditByMitake() {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sms/credit/mitake", nil)
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")
	if res.Code == http.StatusOK {
		var result *model.Result
		json.NewDecoder(res.Body).Decode(&result)

		credit := result.Data.(float64)
		suite.GreaterOrEqual(int(credit), 0, "餘額應大於或等於 0")
	}
}

func (suite *SMSTestSuite) TestSendSMSByEvery8D() {
	b, _ := json.Marshal(suite.sms)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", `/api/v1/sms/every8d/send`, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")
	if res.Code == http.StatusOK {
		var result *model.Result
		var smsResult *model.SMSResult

		json.NewDecoder(res.Body).Decode(&result)
		jsonStr, _ := json.Marshal(result.Data)
		json.Unmarshal(jsonStr, &smsResult)

		suite.GreaterOrEqual(smsResult.Credit, 0, "餘額應大於或等於 0")
	}
}

func (suite *SMSTestSuite) TestSendSMSByMitake() {
	b, _ := json.Marshal(suite.sms)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", `/api/v1/sms/mitake/send`, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")
	if res.Code == http.StatusOK {
		var result *model.Result
		var smsResult *model.SMSResult

		json.NewDecoder(res.Body).Decode(&result)
		jsonStr, _ := json.Marshal(result.Data)
		json.Unmarshal(jsonStr, &smsResult)

		suite.GreaterOrEqual(smsResult.Credit, 0, "餘額應大於或等於 0")
	}
}

func TestSMSTestSuite(t *testing.T) {
	suite.Run(t, new(SMSTestSuite))
}
