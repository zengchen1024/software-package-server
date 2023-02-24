package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type SoftwarePkgBasicInfo struct {
	Id        string
	Importer  dp.Account
	PkgName   dp.PackageName // can't change
	Status    dp.PackageStatus
	AppliedAt int64
}

type SoftwarePkgIssueInfo struct {
	Application Application
	ApprovedBy  []dp.Account
	RejectedBy  []dp.Account
	Comments    []Comment
}

type ImportedSoftwarePkgInfo struct {
	RepoLink dp.URL
}

type SoftwarePkg struct {
	SoftwarePkgBasicInfo

	SoftwarePkgIssueInfo

	ImportedSoftwarePkgInfo
}

func NewSoftwarePkg(user dp.Account, app *Application) SoftwarePkg {
	basic := SoftwarePkgBasicInfo{
		Importer: user,
		PkgName:  app.PackageName,
		Status:   dp.PackageStatusInProgress,
	}

	v := SoftwarePkg{}
	v.SoftwarePkgBasicInfo = basic
	v.Application = *app

	return v
}

// change the status of "creating repo"
// send out the event
// notify the importer
func (entity *SoftwarePkg) Approve(user dp.Account) bool {
	entity.ApprovedBy = append(entity.ApprovedBy, user)
	if entity.IsApproved() {
		entity.Status = dp.PackageStatusCreatingRepo

		return true
	}

	return false
}

// notify the importer
func (entity *SoftwarePkg) Reject() {
	entity.Status = dp.PackageStatusRejected
}

func (entity *SoftwarePkg) GiveUp(user dp.Account) error {
	if user.Account() != entity.Importer.Account() {
		return errors.New("you are not the importer")
	}

	entity.Status = dp.PackageStatusGivingUp

	return nil
}

func (entity *SoftwarePkg) IsApproved() bool {
	return len(entity.ApprovedBy) >= 2
}
