package notification

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"

	"github.com/mirror520/notification/message"
	"github.com/mirror520/notification/model"
)

func SendHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg *message.Message
		if err := ctx.ShouldBind(&msg); err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		_, err := endpoint(ctx, msg)
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
			return
		}

		result := model.SuccessResult("message sent")
		ctx.JSON(http.StatusOK, result)
	}
}

func CreditHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerID := ctx.Param("id")
		if providerID == "" {
			err := errors.New("invalid provider")
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		resp, err := endpoint(ctx, providerID)
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
			return
		}

		result := model.SuccessResult("credit queried")
		result.Data = resp
		ctx.JSON(http.StatusOK, result)
	}
}

func CallbackHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerID := ctx.Param("id")
		if providerID == "" {
			err := errors.New("invalid provider")
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
			return
		}

		req := &CallbackRequest{
			Provider: providerID,
			Values:   ctx.Request.URL.Query(),
		}

		resp, err := endpoint(ctx, req)
		if err != nil {
			result := model.FailureResult(err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, result)
			return
		}

		response, ok := resp.(string)
		if ok {
			ctx.String(http.StatusOK, response)
		} else {
			ctx.JSON(http.StatusOK, resp)
		}
	}
}
