package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/domain/dp"
)

type SoftwarePkgBasicInfo struct {
	Id          string
	Importer    dp.Account
	PackageName dp.PackageName // can't change

	Status     dp.PackageStatus
	ApprovedBy []dp.Account
}

type SoftwarePkg struct {
	SoftwarePkgBasicInfo

	Application Application

	Comments []Comment
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
