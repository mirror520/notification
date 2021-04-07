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
		sms.POST("/:sms_id/send", SendSMSHandler)
		// sms.POST("/master/send", SendSMSHandler)
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

	var providerType provider.SMSProviderType
	smsID := ctx.Param("sms_id")
	switch smsID {
	case "every8d":
		providerType = provider.Every8D

	case "mitake":
		providerType = provider.Mitake
	}

	logger = logger.WithFields(log.Fields{
		"provider": smsID,
	})

	smsProvider, err := provider.SMSProviderCreateFactory(providerType)
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

	logger = logger.WithFields(log.Fields{
		"credit": credit,
	})

	result := model.NewSuccessResult().SetLogger(logger)
	result.AddInfo("查詢餘額成功")
	result.SetData(credit)

	ctx.JSON(http.StatusOK, result)
}

func SendSMSHandler(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"event": "SendSMS",
	})

	var providerType provider.SMSProviderType
	smsID := ctx.Param("sms_id")
	switch smsID {
	case "master":
		providerType = provider.Every8D

	case "every8d":
		providerType = provider.Every8D

	case "mitake":
		providerType = provider.Mitake

	}

	logger = logger.WithFields(log.Fields{
		"provider": smsID,
	})

	smsProvider, err := provider.SMSProviderCreateFactory(providerType)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	var sms model.SMS
	if err := ctx.ShouldBind(&sms); err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo("您輸入的資料格式錯誤")
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	if gin.Mode() != gin.TestMode {
		logger = logger.WithFields(log.Fields{
			"phone": sms.Phone,
		})
	}

	smsResult, err := smsProvider.SendSMS(&sms)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo("簡訊發送失敗")
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	result := model.NewSuccessResult().SetLogger(logger)
	result.AddInfo("簡訊發送成功")
	result.SetData(&smsResult)

	ctx.JSON(http.StatusOK, result)
}

func SwitchSMSMasterHandler(ctx *gin.Context) {
}
