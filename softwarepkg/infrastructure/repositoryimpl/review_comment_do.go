package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	fieldCreatedAt     = "create_at"
	fieldSoftwarePkgId = "software_pkg_id"
)

type SoftwarePkgReviewCommentDO struct {
	Id        uuid.UUID `gorm:"column:id;type:uuid"`
	PkgId     string    `gorm:"column:software_pkg_id;type:uuid"`
	Content   string    `gorm:"column:content"`
	Author    string    `gorm:"column:author"`
	Version   int       `gorm:"column:version"`
	CreatedAt int64     `gorm:"column:created_at"`
	UpdatedAt int64     `gorm:"column:updated_at"`
}

func (s SoftwarePkgReviewCommentDO) toSoftwarePkgReviewComment() (
	r domain.SoftwarePkgReviewComment, err error,
) {
	r.Id = s.Id.String()
	r.CreatedAt = s.CreatedAt

	if r.Author, err = dp.NewAccount(s.Author); err != nil {
		return
	}

	r.Content, err = dp.NewReviewComment(s.Content)

	return
}
