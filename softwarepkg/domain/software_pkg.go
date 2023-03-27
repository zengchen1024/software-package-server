package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

const minNumOfApprover = 2

type User struct {
	Email   dp.Email
	Account dp.Account
}

// SoftwarePkgApplication
type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageDesc       dp.PackageDesc
	PackagePlatform   dp.PackagePlatform
	ImportingPkgSig   dp.ImportingPkgSig
	ReasonToImportPkg dp.ReasonToImportPkg
}

type SoftwarePkgSourceCode struct {
	SpecURL   dp.URL
	SrcRPMURL dp.URL
}

// SoftwarePkgBasicInfo
type SoftwarePkgBasicInfo struct {
	Id          string
	PkgName     dp.PackageName
	Importer    User
	RepoLink    dp.URL
	Phase       dp.PackagePhase
	Frozen      bool
	AppliedAt   int64
	Application SoftwarePkgApplication
	ApprovedBy  []dp.Account
	RejectedBy  []dp.Account
	RelevantPR  dp.URL
}

func (entity *SoftwarePkgBasicInfo) ReviewResult() dp.PackageReviewResult {
	if len(entity.RejectedBy) > 0 {
		return dp.PkgReviewResultRejected
	}

	if len(entity.ApprovedBy) >= minNumOfApprover {
		return dp.PkgReviewResultApproved
	}

	return nil
}

func (entity *SoftwarePkgBasicInfo) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing()
}

// change the status of "creating repo"
// send out the event
// notify the importer
func (entity *SoftwarePkgBasicInfo) ApproveBy(user dp.Account) (bool, error) {
	if !entity.Phase.IsReviewing() || entity.Frozen {
		return false, errors.New("not ready")
	}

	entity.ApprovedBy = append(entity.ApprovedBy, user)

	approved := false
	// only set the result once to avoid duplicate case.
	if len(entity.ApprovedBy) == minNumOfApprover {
		entity.Phase = dp.PackagePhaseCreatingRepo
		approved = true
	}

	return approved, nil
}

// notify the importer
func (entity *SoftwarePkgBasicInfo) RejectBy(user dp.Account) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, errors.New("can't do this")
	}

	entity.RejectedBy = append(entity.RejectedBy, user)

	rejected := false
	// only set the result once to avoid duplicate case.
	if len(entity.RejectedBy) == 1 {
		entity.Phase = dp.PackagePhaseClosed
		rejected = true
	}

	return rejected, nil
}

func (entity *SoftwarePkgBasicInfo) Abandon(user dp.Account) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user, entity.Importer.Account) {
		return errors.New("not the importer")
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandleCI(success bool, pr dp.URL) (bool, error) {
	if entity.RelevantPR != nil {
		return false, errors.New("only handle CI once")
	}

	if entity.Phase.IsClosed() {
		// already closed
		return true, nil
	}

	if !entity.Phase.IsReviewing() {
		return false, errors.New("can't do this")
	}

	entity.RelevantPR = pr

	if success {
		entity.Frozen = false
	} else {
		entity.Phase = dp.PackagePhaseClosed
	}

	return false, nil
}

func (entity *SoftwarePkgBasicInfo) HandlePkgAlreadyExisted() error {
	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

type RepoCreatedInfo struct {
	Platform dp.PackagePlatform
	RepoLink dp.URL
}

func (entity *SoftwarePkgBasicInfo) HandleRepoCreated(info RepoCreatedInfo) error {
	if entity.RepoLink != nil {
		return errors.New("only do once")
	}

	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	if !dp.IsSamePlatform(entity.Application.PackagePlatform, info.Platform) {
		return errors.New("ignore unmached platform")
	}

	entity.RepoLink = info.RepoLink

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandleCodeSaved() error {
	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	entity.Phase = dp.PackagePhaseImported

	return nil
}

// SoftwarePkg
type SoftwarePkg struct {
	SoftwarePkgBasicInfo

	Comments []SoftwarePkgReviewComment
}

func NewSoftwarePkg(user *User, name dp.PackageName, app *SoftwarePkgApplication) SoftwarePkgBasicInfo {
	return SoftwarePkgBasicInfo{
		PkgName:     name,
		Importer:    *user,
		Phase:       dp.PackagePhaseReviewing,
		Frozen:      true,
		Application: *app,
		AppliedAt:   utils.Now(),
	}
}
