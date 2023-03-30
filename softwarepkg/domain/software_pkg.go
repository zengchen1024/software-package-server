package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

// TODO
const minNumOfApprover = 2

type Importer struct {
	Email   dp.Email
	Account dp.Account
}

type User struct {
	Importer

	GiteeID string
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
	Importer    Importer
	RepoLink    dp.URL
	Phase       dp.PackagePhase
	CIStatus    dp.PackageCIStatus
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
func (entity *SoftwarePkgBasicInfo) ApproveBy(user *User) (bool, error) {
	if !entity.Phase.IsReviewing() || entity.Frozen {
		return false, errors.New("not ready")
	}

	entity.ApprovedBy = append(entity.ApprovedBy, user.Account)

	approved := false
	// only set the result once to avoid duplicate case.
	if len(entity.ApprovedBy) == minNumOfApprover {
		entity.Phase = dp.PackagePhaseCreatingRepo
		approved = true
	}

	return approved, nil
}

// notify the importer
func (entity *SoftwarePkgBasicInfo) RejectBy(user *User) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, errors.New("can't do this")
	}

	entity.RejectedBy = append(entity.RejectedBy, user.Account)

	rejected := false
	// only set the result once to avoid duplicate case.
	if len(entity.RejectedBy) == 1 {
		entity.Phase = dp.PackagePhaseClosed
		rejected = true
	}

	return rejected, nil
}

func (entity *SoftwarePkgBasicInfo) Abandon(user *User) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
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
		entity.CIStatus = dp.PackageCIStatusPassed
	} else {
		entity.CIStatus = dp.PackageCIStatusFailed
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

func (entity *SoftwarePkgBasicInfo) UpdateApplication(cmd *SoftwarePkgApplication, user *User) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return errors.New("not the importer")
	}

	entity.Application = *cmd

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
		Importer:    user.Importer,
		Phase:       dp.PackagePhaseReviewing,
		CIStatus:    dp.PackageCIStatusWaiting,
		Frozen:      true,
		Application: *app,
		AppliedAt:   utils.Now(),
	}
}
