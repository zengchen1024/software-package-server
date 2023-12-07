package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequestBody(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		sendFailedResp(ctx, errorBadRequestBody, err)
	}
}

func SendBadRequestParam(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		sendFailedResp(ctx, errorBadRequestParam, err)
	}
}

func SendRespOfCreate(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
}

func SendRespOfPut(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, newResponseCodeMsg("", "success"))
}

func SendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}

func SendRespOfPost(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

func sendFailedResp(ctx *gin.Context, code string, err error) {
	if code == "" {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeMsg(errorSystemError, err.Error()),
		)
	} else {
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(code, err.Error()),
		)
	}
}

func SendError(ctx *gin.Context, err error) {
	sc, code := httpError(err)

	ctx.JSON(sc, newResponseCodeMsg(code, err.Error()))
}
