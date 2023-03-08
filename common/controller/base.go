package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequestBody(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(errorBadRequestBody, err.Error()))
}

func SendBadRequestParam(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(errorBadRequestParam, err.Error()))
}

func SendCreateSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
}

func SendBadRequest(ctx *gin.Context, code string, err error) {
	if code == "" {
		code = errorBadRequest
	}

	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(code, err.Error()))
}

func SendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}
