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
		sms.GET("/credit/:id", SMSCreditHandler)
		sms.POST("/:id/send", SendSMSHandler)
		sms.PATCH("/switch/:id/master", SwitchSMSMasterHandler)
	}
	return router
}

func main() {
	provider.Init()

	router := setRouter()
	router.Run(":7080")
}

func SMSStatusHandler(ctx *gin.Context) {
}

func SMSCreditHandler(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"event": "SMSCredit",
	})

	id := ctx.Param("id")
	p := provider.SMSProviderPool[id]

	logger = logger.WithFields(log.Fields{
		"id":       id,
		"provider": p.Profile().ProviderType(),
	})

	credit, err := p.Credit()
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

	id := ctx.Param("id")
	p := provider.SMSProviderPool[id]

	logger = logger.WithFields(log.Fields{
		"id":       id,
		"provider": p.Profile().ProviderType(),
	})

	var sms model.SMS
	if err := ctx.ShouldBind(&sms); err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo("您輸入的資料格式錯誤")
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	logger = logger.WithFields(log.Fields{
		"phone": sms.Phone,
	})

	smsResult, err := p.SendSMS(&sms)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
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
