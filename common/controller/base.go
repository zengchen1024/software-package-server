package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequestBody(ctx *gin.Context, err error) {
	SendFailedResp(ctx, errorBadRequestBody, err)
}

func SendBadRequestParam(ctx *gin.Context, err error) {
	SendFailedResp(ctx, errorBadRequestParam, err)
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

func SendFailedResp(ctx *gin.Context, code string, err error) {
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
