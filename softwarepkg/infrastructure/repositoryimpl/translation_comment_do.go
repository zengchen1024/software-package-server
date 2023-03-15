package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

type SoftwarePkgTranslationCommentDO struct {
	// must set "uuid" as the name of column
	Id        uuid.UUID `gorm:"column:uuid;type:uuid"`
	PkgId     string    `gorm:"column:software_pkg_id"`
	Content   string    `gorm:"column:content"`
	Language  string    `gorm:"column:language"`
	CommentId string    `gorm:"column:review_comment_id"`
	CreatedAt int64     `gorm:"column:created_at"`
	UpdatedAt int64     `gorm:"column:updated_at"`
	Version   int       `gorm:"column:version"`
}

func (t translationComment) toSoftwarePkgTranslationCommentDO(
	pid string, comment *domain.SoftwarePkgTranslatedReviewComment, do *SoftwarePkgTranslationCommentDO,
) {
	*do = SoftwarePkgTranslationCommentDO{
		Id:        uuid.New(),
		PkgId:     pid,
		Content:   comment.Content,
		Language:  comment.Language.Language(),
		CommentId: comment.CommentId,
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	}
}

func (do *SoftwarePkgTranslationCommentDO) toSoftwarePkgTranslatedReviewComment() (
	r domain.SoftwarePkgTranslatedReviewComment, err error,
) {
	r.Id = do.Id.String()
	r.CommentId = do.CommentId
	r.Content = do.Content

	r.Language, err = dp.NewLanguage(do.Language)

	return
}
