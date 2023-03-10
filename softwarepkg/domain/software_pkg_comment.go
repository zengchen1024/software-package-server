package domain

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

// SoftwarePkgReviewComment
type SoftwarePkgReviewComment struct {
	Id        string
	CreatedAt int64
	Author    dp.Account
	Content   dp.ReviewComment
	Language  string
}

func NewSoftwarePkgReviewComment(
	author dp.Account, content dp.ReviewComment,
) SoftwarePkgReviewComment {
	return SoftwarePkgReviewComment{
		CreatedAt: utils.Now(),
		Author:    author,
		Content:   content,
	}
}

type SoftwarePkgTranslatedReviewComment struct {
	Id        string
	CommentId string
	Language  string
	Content   dp.ReviewComment
}
