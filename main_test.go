package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/notification/model"
	"github.com/mirror520/notification/provider"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type SMSCreditTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *SMSCreditTestSuite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)

	provider.Init()
	suite.router = internalRouter()
}

func (suite *SMSCreditTestSuite) TestSMSCreditByEvery8D() {
	SMSCredit(suite, "every8d")
}

func (suite *SMSCreditTestSuite) TestSMSCreditByMitake() {
	SMSCredit(suite, "mitake")
}

func SMSCredit(suite *SMSCreditTestSuite, pid string) {
	resource := fmt.Sprintf("/api/v1/sms/credit/%s", pid)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", resource, nil)
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

type SMSSendTestSuite struct {
	suite.Suite
	router *gin.Engine
	sms    model.SMS
}

func (suite *SMSSendTestSuite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)

	provider.Init()
	suite.router = internalRouter()
	suite.sms = model.SMS{
		Phone:   os.Getenv("SMS_TESTPHONE"),
		Message: "現在時間: " + time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (suite *SMSSendTestSuite) TestSendSMSByEvery8D() {
	sms := suite.sms
	sms.Message += ", 簡訊服務商: Every8D"
	SendSMS(suite, sms, "every8d")
}

func (suite *SMSSendTestSuite) TestSendSMSByMitake() {
	sms := suite.sms
	sms.Message += ", 簡訊服務商: Mitake"
	SendSMS(suite, sms, "mitake")
}

func (suite *SMSSendTestSuite) TestSendSMSByMaster() {
	sms := suite.sms
	sms.Message += ", 簡訊服務商: Every8D (master)"
	SendSMS(suite, sms)

	masterID := provider.SMSMasterProvider().Profile().ID
	suite.Equal("every8d", masterID, "預設主要簡訊提供商為 master")
}

func SendSMS(suite *SMSSendTestSuite, sms model.SMS, pid ...string) {
	b, _ := json.Marshal(sms)

	resource := "/api/v1/sms/send"
	if len(pid) > 0 {
		resource += fmt.Sprintf("/%s", pid[0])
	}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", resource, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	suite.router.ServeHTTP(res, req)

	var result *model.Result
	if sms.Phone != "" {
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

type SMSSwitchMasterTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *SMSSwitchMasterTestSuite) SetupSuite() {
	logrus.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)

	provider.Init()
	suite.router = internalRouter()
}

func (suite *SMSSwitchMasterTestSuite) TestSMSSwitchMaster() {
	resource := fmt.Sprintf("/api/v1/sms/switch/%s/master", "mitake")

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", resource, nil)
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "HTTP 狀態碼應為 200")
	if res.Code == http.StatusOK {
		var result *model.Result
		json.NewDecoder(res.Body).Decode(&result)

		suite.Equal(result.Status, model.Success, "狀態應為成功")

		masterID := provider.SMSMasterProvider().Profile().ID
		suite.Equal("mitake", masterID, "切換主要簡訊提供商為 Mitake")
	}
}

func TestSMSCreditTestSuite(t *testing.T) {
	suite.Run(t, new(SMSCreditTestSuite))
}

func TestSMSSendTestSuite(t *testing.T) {
	suite.Run(t, new(SMSSendTestSuite))
}

func TestSMSSwitchMasterTestSuite(t *testing.T) {
	suite.Run(t, new(SMSSwitchMasterTestSuite))
}
