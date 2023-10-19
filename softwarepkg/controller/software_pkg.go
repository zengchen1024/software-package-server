package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	commonctl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type SoftwarePkgController struct {
	service app.SoftwarePkgService
}

func AddRouteForSoftwarePkgController(r *gin.RouterGroup, pkgService app.SoftwarePkgService) {
	ctl := SoftwarePkgController{
		service: pkgService,
	}

	m := middleware.UserChecking().CheckUser
	r.POST("/v1/softwarepkg", m, ctl.ApplyNewPkg)
	r.GET("/v1/softwarepkg", ctl.ListPkgs)
	r.GET("/v1/softwarepkg/:id", ctl.Get)
	r.PUT("/v1/softwarepkg/:id", m, ctl.UpdateApplication)

	r.POST("/v1/softwarepkg/:id/review", m, ctl.Review)
	r.PUT("/v1/softwarepkg/:id/reject", m, ctl.Reject)
	r.PUT("/v1/softwarepkg/:id/abandon", m, ctl.Abandon)
	r.PUT("/v1/softwarepkg/:id/review/rerunci", m, ctl.RerunCI)
	r.POST("/v1/softwarepkg/:id/review/comment", m, ctl.NewReviewComment)
	r.POST("/v1/softwarepkg/:id/review/comment/:cid/translate", m, ctl.TranslateReviewComment)
}

// ApplyNewPkg
// @Summary apply a new software package
// @Description apply a new software package
// @Tags  SoftwarePkg
// @Accept json
// @Param	param  body	 softwarePkgRequest	 true	"body of applying a new software package"
// @Success 201 {object} app.NewSoftwarePkgDTO
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg [post]
func (ctl SoftwarePkgController) ApplyNewPkg(ctx *gin.Context) {
	var req softwarePkgRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	cmd, err := req.toCmd(&user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if r, code, err := ctl.service.ApplyNewPkg(&cmd); err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfPost(ctx, r)
	}
}

// ListPkgs
// @Summary list software packages
// @Description list software packages
// @Tags  SoftwarePkg
// @Accept json
// @Param    importer         query	 string   false    "importer of the softwarePkg"
// @Param    phase            query	 string   false    "phase of the softwarePkg"
// @Param    count_per_page   query	 int      false    "count per page"
// @Param    page_num         query	 int      false    "page num which starts from 1"
// @Success 200 {object} app.SoftwarePkgsDTO
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg [get]
func (ctl SoftwarePkgController) ListPkgs(ctx *gin.Context) {
	var req softwarePkgListQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if v, err := ctl.service.ListPkgs(&cmd); err != nil {
		commonctl.SendFailedResp(ctx, "", err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// Get
// @Summary get software package
// @Description get software package
// @Tags  SoftwarePkg
// @Accept json
// @Param    id         path	string  true    "id of software package"
// @Success 200 {object} app.SoftwarePkgReviewDTO
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id} [get]
func (ctl SoftwarePkgController) Get(ctx *gin.Context) {
	if v, code, err := ctl.service.GetPkgReviewDetail(ctx.Param("id")); err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// Review
// @Summary review software package
// @Description review software package
// @Tags  SoftwarePkg
// @Accept json
// @Param  id     path   string         true  "id of software package"
// @Param  param  body   reviewRequest  true  "body of reviewing a software package"
// @Success 201 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/review [post]
func (ctl SoftwarePkgController) Review(ctx *gin.Context) {
	var req reviewRequest

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	info, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	if err := ctl.service.Review(ctx.Param("id"), &user, info); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}

// Reject
// @Summary reject software package
// @Description reject software package
// @Tags  SoftwarePkg
// @Accept json
// @Param	id  path	 string	 true	"id of software package"
// @Success 202 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/reject [put]
func (ctl SoftwarePkgController) Reject(ctx *gin.Context) {
	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	if err := ctl.service.Reject(ctx.Param("id"), &user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}

// Abandon
// @Summary abandon software package
// @Description abandon software package
// @Tags  SoftwarePkg
// @Accept json
// @Param	id  path	 string	 true	"id of software package"
// @Success 202 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/abandon [put]
func (ctl SoftwarePkgController) Abandon(ctx *gin.Context) {
	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	if err := ctl.service.Abandon(ctx.Param("id"), &user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}

// NewReviewComment
// @Summary create a new software package review comment
// @Description create a new software package review comment
// @Tags  SoftwarePkg
// @Accept json
// @Param	param  body	 reviewCommentRequest	 true	"body of creating a new software package review comment"
// @Param	id     path	 string	                 true	"id of software package"
// @Success 201 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/review/comment [post]
func (ctl SoftwarePkgController) NewReviewComment(ctx *gin.Context) {
	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	var req reviewCommentRequest
	if err = ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(&user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if code, err := ctl.service.NewReviewComment(ctx.Param("id"), &cmd); err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfCreate(ctx)
	}
}

// TranslateReviewComment
// @Summary translate review comment
// @Description translate review comment
// @Tags  SoftwarePkg
// @Accept json
// @Param    id       path       string                      true    "id of software package"
// @Param    cid      path       string                      true    "cid of review comment"
// @Param    param    body       translationCommentRequest   true    "body of translate review comment"
// @Success 201 {object} app.TranslatedReveiwCommentDTO
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/review/comment/{cid}/translate [post]
func (ctl SoftwarePkgController) TranslateReviewComment(ctx *gin.Context) {
	var req translationCommentRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("id"), ctx.Param("cid"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if v, code, err := ctl.service.TranslateReviewComment(&cmd); err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// UpdateApplication
// @Summary update application of software package
// @Description update application of software package
// @Tags  SoftwarePkg
// @Accept json
// @Param    id      path     string                true    "id of software package"
// @Param    param   body     softwarePkgRequest    true    "body of updating software package application"
// @Success 202 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/:id [put]
func (ctl SoftwarePkgController) UpdateApplication(ctx *gin.Context) {
	var req softwarePkgRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	cmd, err := req.toCmd(&user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	code, err := ctl.service.UpdateApplication(
		&app.CmdToUpdateSoftwarePkgApplication{
			PkgId:                    ctx.Param("id"),
			CmdToApplyNewSoftwarePkg: cmd,
		},
	)
	if err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}

// ReRunCI
// @Summary rerun ci of software package
// @Description rerun ci of software package
// @Tags  SoftwarePkg
// @Accept json
// @Param	id  path	 string	 true	"id of software package"
// @Success 202 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/review/rerunci [put]
func (ctl SoftwarePkgController) RerunCI(ctx *gin.Context) {
	user, err := middleware.UserChecking().FetchUser(ctx)
	if err != nil {
		commonctl.SendFailedResp(ctx, "", err)

		return
	}

	if code, err := ctl.service.RerunCI(ctx.Param("id"), &user); err != nil {
		commonctl.SendFailedResp(ctx, code, err)
	} else {
		commonctl.SendRespOfPut(ctx)
	}
}
