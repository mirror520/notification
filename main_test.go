package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mirror520/sms/model"
	"github.com/stretchr/testify/assert"
)

func TestSMSCredit(t *testing.T) {
	router := setRouter()

	res := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/sms/credit/every8d", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "成功取得餘額")

	var result *model.Result
	json.NewDecoder(res.Body).Decode(&result)

	credit := result.Data.(float64)
	assert.LessOrEqual(t, 0.0, credit)
}
