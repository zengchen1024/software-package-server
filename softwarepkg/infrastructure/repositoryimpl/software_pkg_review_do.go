package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	softwareUuid = "software_pkg_uuid"
	createAt     = "create_time"
)

type SoftwarePkgReviewDO struct {
	UUID            uuid.UUID `gorm:"column:uuid;type:uuid"`
	Content         string    `gorm:"column:content"`
	ApplyUser       string    `gorm:"column:apply_user"`
	SoftwarePkgUUID string    `gorm:"column:software_pkg_uuid;type:uuid"`
	Status          int       `gorm:"column:status"`
	Version         int       `gorm:"column:version"`
	CreatedAt       int64     `gorm:"column:create_time"`
	UpdatedAt       int64     `gorm:"column:update_time"`
}

func (s SoftwarePkgReviewDO) toSoftwarePkgReviewCommentSummary() (pkgComment domain.SoftwarePkgReviewComment, err error) {
	pkgComment.CreatedAt = s.CreatedAt

	pkgComment.Id = s.UUID.String()

	if pkgComment.Author, err = dp.NewAccount(s.ApplyUser); err != nil {
		return
	}

	pkgComment.Content, err = dp.NewReviewComment(s.Content)

	return
}
