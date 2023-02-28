package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	commonctl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type SoftwarePkgController struct {
	commonctl.BaseController
	service app.SoftwarePkgService
}

func AddRouteForSoftwarePkgController(r *gin.RouterGroup, pkgService app.SoftwarePkgService) {
	ctl := SoftwarePkgController{
		service: pkgService,
	}

	r.POST("/v1/softwarepkg", ctl.ApplyNewPkg)

}

// ApplyNewPkg
// @Summary apply a new software package
// @Description apply a new software package
// @Tags  SoftwarePkg
// @Accept json
// @Param	param  body	 softwarePkgRequest	 true	"body of applying a new software package"
// @Success 201 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg [post]
func (ctl SoftwarePkgController) ApplyNewPkg(ctx *gin.Context) {
	var req softwarePkgRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ctl.SendBadRequestBody(ctx, err)

		return
	}

	pkg, err := req.toCmd()
	if err != nil {
		ctl.SendBadRequestParam(ctx, err)

		return
	}

	if code, err := ctl.service.ApplyNewPkg(nil, &pkg); err != nil {
		ctl.SendBadRequest(ctx, code, err)
	} else {
		ctl.SendCreateSuccess(ctx)
	}
}
