package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	fieldCreatedAt     = "created_at"
	fieldSoftwarePkgId = "software_pkg_id"
)

type SoftwarePkgReviewCommentDO struct {
	Id        uuid.UUID `gorm:"column:id;type:uuid"`
	PkgId     string    `gorm:"column:software_pkg_id;type:uuid"`
	Content   string    `gorm:"column:content"`
	Author    string    `gorm:"column:author"`
	CreatedAt int64     `gorm:"column:created_at"`
	UpdatedAt int64     `gorm:"column:updated_at"`
	Version   int       `gorm:"column:version"`
}

func (do *SoftwarePkgReviewCommentDO) toSoftwarePkgReviewComment() (
	r domain.SoftwarePkgReviewComment, err error,
) {
	r.Id = do.Id.String()
	r.CreatedAt = do.CreatedAt

	if r.Author, err = dp.NewAccount(do.Author); err != nil {
		return
	}

	r.Content, err = dp.NewReviewComment(do.Content)

	return
}
