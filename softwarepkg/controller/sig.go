package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/sigvalidator"
)

type SigController struct {
	service sigvalidator.SigValidator
}

func AddRouteForSigController(r *gin.RouterGroup, sigService sigvalidator.SigValidator) {
	ctl := SigController{
		service: sigService,
	}

	r.GET("/v1/sig", middleware.UserChecking().CheckUser, ctl.List)
}

// List
// @Summary list sigs
// @Description list sigs
// @Tags  Sig
// @Accept json
// @Success 200 {object} sigvalidator.Sig
// @Router /v1/sig [get]
func (s SigController) List(ctx *gin.Context) {
	commonctl.SendRespOfGet(ctx, s.service.GetAll())
}
