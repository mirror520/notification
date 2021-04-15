package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirror520/sms/model"
	"github.com/mirror520/sms/provider"

	log "github.com/sirupsen/logrus"
)

func internalRouter() *gin.Engine {
	router := gin.Default()

	sms := router.Group("/api/v1/sms")
	{
		sms.GET("/status", SMSStatusHandler)
		sms.GET("/credit/:pid", SMSCreditHandler)
		sms.POST("/send", SendSMSHandler)
		sms.POST("/send/:pid", SendSMSHandler)
		sms.PATCH("/switch/:pid/master", SwitchSMSMasterHandler)
	}
	return router
}

func externalRouter() *gin.Engine {
	router := gin.Default()

	sms := router.Group("/api/v1/sms")
	{
		sms.GET("/status/:pid/callback", SMSStatusCallbackHandler)
	}

	return router
}

func main() {
	provider.Init()

	go internalRouter().Run(":7080")
	externalRouter().Run(":7090")
}

func SMSStatusHandler(ctx *gin.Context) {
}

func SMSCreditHandler(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"event": "SMSCredit",
	})

	pid := ctx.Param("pid")
	p, err := provider.SMSProvider(pid)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	logger = logger.WithFields(log.Fields{
		"pid":      pid,
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

	pid := ctx.Param("pid")
	p, err := provider.SMSProvider(pid)

	// GET /api/v1/sms/send (master)
	if pid == "" {
		p = provider.SMSMasterProvider()
		pid = p.Profile().ID
	}

	// GET /api/v1/sms/send/:pid
	if err != nil && p == nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	logger = logger.WithFields(log.Fields{
		"pid":      pid,
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
	logger := log.WithFields(log.Fields{
		"event": "SwitchSMSMaster",
	})

	pid := ctx.Param("pid")
	p, err := provider.SMSProvider(pid)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	provider.SwitchSMSProviderToMaster(p)

	result := model.NewSuccessResult().SetLogger(logger)
	result.AddInfo("切換主要簡訊提供商成功")

	ctx.JSON(http.StatusOK, result)
}

func SMSStatusCallbackHandler(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"event": "SMSStatusCallback",
	})

	pid := ctx.Param("pid")
	p, err := provider.SMSProvider(pid)
	if err != nil {
		result := model.NewFailureResult().SetLogger(logger)
		result.AddInfo(err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
		return
	}

	queryParams := ctx.Request.URL.Query()
	phone, response := p.Callback(&queryParams)

	logger = logger.WithFields(log.Fields{
		"pid":      pid,
		"provider": p.Profile().ProviderType(),
		"phone":    phone,
	})

	logger.Infoln("成功接收簡訊狀態")
	ctx.String(http.StatusOK, response)
}
