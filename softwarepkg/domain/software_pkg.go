package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type SoftwarePkgReviewComment struct {
	Id        string
	CreatedAt int64
	Author    dp.Account
	Content   dp.ReviewComment
}

// SoftwarePkgApplication
type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageName       dp.PackageName
	PackageDesc       dp.PackageDesc
	PackagePlatform   dp.PackagePlatform
	ImportingPkgSig   dp.ImportingPkgSig
	ReasonToImportPkg dp.ReasonToImportPkg
}

type SoftwarePkgSourceCode struct {
	Address dp.URL
	License dp.License
}

// SoftwarePkgBasicInfo
type SoftwarePkgBasicInfo struct {
	Id           string
	PkgName      dp.PackageName // can't change
	Importer     dp.Account
	RepoLink     dp.URL
	Phase        dp.PackagePhase
	ReviewResult dp.PackageReviewResult
	AppliedAt    int64
	Application  SoftwarePkgApplication
	ApprovedBy   []dp.Account
	RejectedBy   []dp.Account
}

func (entity *SoftwarePkgBasicInfo) removeFromRejected(u dp.Account) bool {
	i := -1
	for j := range entity.RejectedBy {
		if dp.IsSameAccount(entity.RejectedBy[j], u) {
			i = j

			break
		}
	}

	if i < 0 {
		return false
	}

	n := len(entity.RejectedBy) - 1

	if i == 0 {
		if n == 0 {
			entity.RejectedBy = nil
		} else {
			entity.RejectedBy = entity.RejectedBy[1:]
		}
	} else {
		if i != n {
			entity.RejectedBy[i] = entity.RejectedBy[n]
		}
		entity.RejectedBy = entity.RejectedBy[:n]
	}

	return true
}

// change the status of "creating repo"
// send out the event
// notify the importer
func (entity *SoftwarePkgBasicInfo) ApproveBy(user dp.Account) (changed, approved bool) {
	if !dp.IsPackagePhaseReviewing(entity.Phase) {
		return
	}

	entity.ApprovedBy = append(entity.ApprovedBy, user)
	changed = true

	if len(entity.RejectedBy) > 0 {
		if !entity.removeFromRejected(user) || len(entity.RejectedBy) != 0 {
			return
		}

		if len(entity.ApprovedBy) >= 2 {
			entity.ReviewResult = dp.PkgReviewResultApproved
			approved = true
		}
	} else {
		// only set the result once to avoid duplicate case.
		if len(entity.ApprovedBy) == 2 {
			entity.ReviewResult = dp.PkgReviewResultApproved
			approved = true
		}
	}

	return
}

// notify the importer
func (entity *SoftwarePkgBasicInfo) RejectBy(user dp.Account) (changed, rejected bool) {
	if !dp.IsPackagePhaseReviewing(entity.Phase) {
		return
	}

	entity.RejectedBy = append(entity.RejectedBy, user)
	changed = true

	// only set the result once to avoid duplicate case.
	if len(entity.RejectedBy) == 1 {
		entity.ReviewResult = dp.PkgReviewResultRejected
		rejected = true
	}

	return
}

func (entity *SoftwarePkgBasicInfo) GiveUp(user dp.Account) error {
	if !dp.IsPackagePhaseReviewing(entity.Phase) {
		return errors.New("can't do this")
	}

	if user.Account() != entity.Importer.Account() {
		return errors.New("you are not the importer")
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

func (entity *SoftwarePkgBasicInfo) Close() error {
	if !dp.IsPackagePhaseReviewing(entity.Phase) {
		return errors.New("can't do this")
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

// SoftwarePkg
type SoftwarePkg struct {
	SoftwarePkgBasicInfo

	Comments []SoftwarePkgReviewComment
}

func NewSoftwarePkg(user dp.Account, app *SoftwarePkgApplication) SoftwarePkg {
	basic := SoftwarePkgBasicInfo{
		Importer: user,
		PkgName:  app.PackageName,
		Phase:    dp.PackagePhaseReviewing,
	}

	v := SoftwarePkg{}
	v.SoftwarePkgBasicInfo = basic
	v.Application = *app

	return v
}
