package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/clavalidator"
)

type CLAController struct {
	service clavalidator.ClaValidator
}

func AddRouteForCLAController(r *gin.RouterGroup, service clavalidator.ClaValidator) {
	ctl := CLAController{
		service: service,
	}

	r.GET("/v1/cla", middleware.UserChecking().CheckUser, ctl.VerifyCla)
}

// VerifyCla
// @Summary verify cla
// @Description verify cla
// @Tags  CLA
// @Accept json
// @Success 200 {object} claSingedResp
// @Failure 400 {object} ResponseData
// @Router /v1/cla [get]
func (c CLAController) VerifyCla(ctx *gin.Context) {
	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := c.service.HasSignedCLA(user.Email); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, claSingedResp{v})
	}
}

type claSingedResp struct {
	Signed bool `json:"signed"`
}
