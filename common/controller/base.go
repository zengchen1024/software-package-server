package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (ctl BaseController) SendBadRequestBody(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(errorBadRequestBody, err.Error()))
}

func (ctl BaseController) SendBadRequestParam(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(errorBadRequestParam, err.Error()))
}

func (ctl BaseController) SendCreateSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
}

func (ctl BaseController) SendBadRequest(ctx *gin.Context, code string, err error) {
	if code == "" {
		code = errorBadRequest
	}

	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(code, err.Error()))
}

func (ctl BaseController) SendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}
