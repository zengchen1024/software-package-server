package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

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

// SoftwarePkgCI
type SoftwarePkgCI struct {
	PRNum  int
	Status dp.PackageCIStatus
	// TODO deal with the case that the ci is timeout
	//startTime int64
}

func (ci *SoftwarePkgCI) isSuccess() bool {
	return ci.Status != nil && ci.Status.IsCIPassed()
}

// SoftwarePkgBasicInfo
type SoftwarePkgBasicInfo struct {
	Id          string
	PkgName     dp.PackageName
	Importer    Importer
	RepoLink    dp.URL
	Phase       dp.PackagePhase
	CI          SoftwarePkgCI
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

	if len(entity.ApprovedBy) >= config.MinNumOfApprovers {
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
	if !entity.Phase.IsReviewing() || !entity.CI.isSuccess() {
		return false, errors.New("can't do this")
	}

	entity.ApprovedBy = append(entity.ApprovedBy, user.Account)

	approved := false
	// only set the result once to avoid duplicate case.
	if len(entity.ApprovedBy) == config.MinNumOfApprovers {
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
		return errorNotTheImporter
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

func (entity *SoftwarePkgBasicInfo) RerunCI(user *User) error {
	b := entity.Phase.IsReviewing() && !entity.CI.Status.IsCIRunning()
	if !b {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Importer.Account, entity.Importer.Account) {
		return errorNotTheImporter
	}

	entity.CI = SoftwarePkgCI{Status: dp.PackageCIStatusWaiting}

	return nil
}

func (entity *SoftwarePkgBasicInfo) UpdateApplication(cmd *SoftwarePkgApplication, user *User) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return errorNotTheImporter
	}

	entity.Application = *cmd

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandleCIChecking() error {
	b := entity.Phase.IsReviewing() && entity.CI.Status.IsCIWaiting()
	if !b {
		return errors.New("can't do this")
	}

	entity.CI.Status = dp.PackageCIStatusRunning
	//entity.CI.startTime = utils.Now()

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandleCIChecked(success bool, prNum int) error {
	if !entity.Phase.IsReviewing() || entity.CI.PRNum != prNum {
		return errors.New("can't do this")
	}

	if success {
		entity.CI.PRNum = prNum
		entity.CI.Status = dp.PackageCIStatusPassed
	} else {
		entity.CI.Status = dp.PackageCIStatusFailed
	}

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandlePkgInitialized(pr dp.URL) error {
	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	entity.RelevantPR = pr

	return nil
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
	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	if !dp.IsSamePlatform(entity.Application.PackagePlatform, info.Platform) {
		return errors.New("ignore unmached platform")
	}

	entity.RepoLink = info.RepoLink

	return nil
}

func (entity *SoftwarePkgBasicInfo) HandleCodeSaved(info RepoCreatedInfo) error {
	if err := entity.HandleRepoCreated(info); err != nil {
		return err
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
		Importer:    user.Importer,
		Phase:       dp.PackagePhaseReviewing,
		CI:          SoftwarePkgCI{Status: dp.PackageCIStatusWaiting},
		Application: *app,
		AppliedAt:   utils.Now(),
	}
}
