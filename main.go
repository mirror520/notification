package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/sms/model"
	"github.com/mirror520/sms/provider"

	log "github.com/sirupsen/logrus"
)

func setRouter() *gin.Engine {
	router := gin.Default()

	sms := router.Group("/api/v1/sms")
	{
		sms.GET("/status", SMSStatusHandler)
		sms.GET("/credit/:sms_id", SMSCreditHandler)
		sms.POST("/:phone/send", SendSMSToPhoneHandler)
		sms.PATCH("/switch/:sms_id/master", SwitchSMSMasterHandler)
	}
	return router
}

func main() {
	router := setRouter()
	router.Run(":7080")
}

func SMSStatusHandler(ctx *gin.Context) {
}

func SMSCreditHandler(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"event": "SMSCredit",
	})

	providerType := provider.Every8D
	smsProvider, err := provider.CreateSMSProviderFactory(providerType)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	credit, err := smsProvider.Credit()
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	result := model.NewSuccessResult().SetLogger(logger)
	result.AddInfo("查詢餘額成功")
	result.SetData(credit)

	ctx.JSON(http.StatusOK, result)
}

func SendSMSToPhoneHandler(ctx *gin.Context) {
}

func SwitchSMSMasterHandler(ctx *gin.Context) {
}
