package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/sms/model"
	"github.com/stretchr/testify/suite"
)

type SMSTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *SMSTestSuite) SetupTest() {
	gin.SetMode(gin.ReleaseMode)
	suite.router = setRouter()
}

func (suite *SMSTestSuite) TestSMSEvery8DCredit() {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sms/credit/every8d", nil)
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "成功取得餘額")
	if res.Code == http.StatusOK {
		var result *model.Result
		json.NewDecoder(res.Body).Decode(&result)

		credit := result.Data.(float64)
		suite.LessOrEqual(0.0, credit, "餘額大於 0")
	}
}

func (suite *SMSTestSuite) TestSMSMitakeCredit() {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sms/credit/mitake", nil)
	suite.router.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code, "成功取得餘額")
	if res.Code == http.StatusOK {
		var result *model.Result
		json.NewDecoder(res.Body).Decode(&result)

		credit := result.Data.(float64)
		suite.LessOrEqual(0.0, credit, "餘額大於 0")
	}
}

func TestSMSTestSuite(t *testing.T) {
	suite.Run(t, new(SMSTestSuite))
}
