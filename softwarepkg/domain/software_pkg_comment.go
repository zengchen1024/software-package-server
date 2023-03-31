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

// SoftwarePkgTranslatedReviewComment
type SoftwarePkgTranslatedReviewComment struct {
	Id        string
	CommentId string
	Content   string
	Language  dp.Language
}

func NewSoftwarePkgTranslatedReviewComment(
	comment *SoftwarePkgReviewComment, content string, lang dp.Language,
) SoftwarePkgTranslatedReviewComment {
	return SoftwarePkgTranslatedReviewComment{
		Content:   content,
		Language:  lang,
		CommentId: comment.Id,
	}
}
