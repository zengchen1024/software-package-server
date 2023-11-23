package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	commonctl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type SoftwarePkgCommentController struct {
	commentApp app.SoftwarePkgCommentAppService
}

func AddRouteForSoftwarePkgCommentController(
	r *gin.RouterGroup,
	commentApp app.SoftwarePkgCommentAppService,
) {
	ctl := SoftwarePkgCommentController{
		commentApp: commentApp,
	}

	m := middleware.UserChecking().CheckUser

	r.GET("/v1/softwarepkg/:id/review/comment", ctl.List)
	r.POST("/v1/softwarepkg/:id/review/comment", m, ctl.NewReviewComment)
	r.POST("/v1/softwarepkg/:id/review/comment/:cid/translate", ctl.TranslateReviewComment)
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
func (ctl SoftwarePkgCommentController) NewReviewComment(ctx *gin.Context) {
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

	cmd, err := req.toCmd(ctx.Param("id"), &user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.commentApp.NewReviewComment(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfCreate(ctx)
	}
}

// ListComment
// @Summary list software package review comment
// @Description list software package review comment
// @Tags  SoftwarePkg
// @Accept json
// @Param	id     path	 string	                 true	"id of software package"
// @Success 202 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/softwarepkg/{id}/review/comment [get]
func (ctl SoftwarePkgCommentController) List(ctx *gin.Context) {
	if v, err := ctl.commentApp.ListComments(ctx.Param("id")); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
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
func (ctl SoftwarePkgCommentController) TranslateReviewComment(ctx *gin.Context) {
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

	if v, err := ctl.commentApp.TranslateReviewComment(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}
