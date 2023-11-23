package controller

import (
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type reviewCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

func (r reviewCommentRequest) toCmd(pkgId string, user *domain.User) (rc app.CmdToWriteSoftwarePkgReviewComment, err error) {
	if rc.Content, err = dp.NewReviewComment(r.Comment); err != nil {
		return
	}

	rc.PkgId = pkgId
	rc.Author = user.Account

	return
}

type translationCommentRequest struct {
	Language string `json:"language"`
}

func (t translationCommentRequest) toCmd(pkgId, commentId string) (cmd app.CmdToTranslateReviewComment, err error) {
	cmd.PkgId = pkgId
	cmd.CommentId = commentId
	cmd.Language, err = dp.NewLanguage(t.Language)

	return
}
