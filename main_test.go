package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/sms/model"
	"github.com/mirror520/sms/provider"
	"github.com/stretchr/testify/suite"

	"github.com/sirupsen/logrus"
)

type SMSTestSuite struct {
	suite.Suite
	router *gin.Engine
	sms    *model.SMS
}

func (suite *SMSTestSuite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)

	provider.Init()
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

		suite.Equal(result.Status, model.Success, "狀態應為成功")
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

		suite.Equal(result.Status, model.Success, "狀態應為成功")
		suite.GreaterOrEqual(int(credit), 0, "餘額應大於或等於 0")
	}
}

func (suite *SMSTestSuite) TestSendSMSByEvery8D() {
	suite.sms.Message += ", 簡訊服務商: Every8D"
	b, _ := json.Marshal(suite.sms)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", `/api/v1/sms/every8d/send`, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	suite.router.ServeHTTP(res, req)

	var result *model.Result
	if suite.sms.Phone != "" {
		suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")

		json.NewDecoder(res.Body).Decode(&result)
		jsonStr, _ := json.Marshal(result.Data)

		var smsResult *model.SMSResult
		json.Unmarshal(jsonStr, &smsResult)

		suite.Equal(result.Status, model.Success, "狀態應為成功")
		suite.GreaterOrEqual(smsResult.Credit, 0, "餘額應大於或等於 0")
	} else {
		suite.Equal(http.StatusUnprocessableEntity, res.Code, "HTTP 狀態碼應為 422")

		json.NewDecoder(res.Body).Decode(&result)
		err := errors.New(result.Info[0])

		suite.Equal(result.Status, model.Failure, "狀態應為失敗")
		suite.Error(err, "應該發生格式錯誤")
	}
}

func (suite *SMSTestSuite) TestSendSMSByMitake() {
	suite.sms.Message += ", 簡訊服務商: Mitake"
	b, _ := json.Marshal(suite.sms)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", `/api/v1/sms/mitake/send`, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	suite.router.ServeHTTP(res, req)

	var result *model.Result
	if suite.sms.Phone != "" {
		suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")

		json.NewDecoder(res.Body).Decode(&result)
		jsonStr, _ := json.Marshal(result.Data)

		var smsResult *model.SMSResult
		json.Unmarshal(jsonStr, &smsResult)

		suite.Equal(result.Status, model.Success, "狀態應為成功")
		suite.GreaterOrEqual(smsResult.Credit, 0, "餘額應大於或等於 0")
	} else {
		suite.Equal(http.StatusUnprocessableEntity, res.Code, "HTTP 狀態碼應為 422")

		json.NewDecoder(res.Body).Decode(&result)
		err := errors.New(result.Info[0])

		suite.Equal(result.Status, model.Failure, "狀態應為失敗")
		suite.Error(err, "應該發生格式錯誤")
	}
}

func TestSMSTestSuite(t *testing.T) {
	suite.Run(t, new(SMSTestSuite))
}
