package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/opensourceways/software-package-server/common/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

var (
	notImporter    = allerror.New(allerror.ErrorCodeNotImporter, "not the importer")
	incorrectPhase = allerror.New(allerror.ErrorCodeIncorrectPhase, "incorrect phase")
)

type User struct {
	Email dp.Email
	// TODO name
	Account dp.Account

	GiteeID  string
	GithubID string
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
	Upstream  dp.URL
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

// SoftwarePkgApprover
type SoftwarePkgApprover struct {
	Account dp.Account
	IsTC    bool
}

func (approver *SoftwarePkgApprover) String() string {
	return fmt.Sprintf("%s/%v", approver.Account.Account(), approver.IsTC)
}

func StringToSoftwarePkgApprover(s string) (r SoftwarePkgApprover, err error) {
	items := strings.Split(s, "/")

	if r.Account, err = dp.NewAccount(items[0]); err == nil {
		r.IsTC, _ = strconv.ParseBool(items[1])
	}

	return
}

// SoftwarePkgBasicInfo
type SoftwarePkgBasicInfo struct {
	Id          string
	PkgName     dp.PackageName
	Importer    User
	RepoLink    dp.URL
	Phase       dp.PackagePhase
	CI          SoftwarePkgCI
	AppliedAt   int64
	Application SoftwarePkgApplication
	RelevantPR  dp.URL
	Review      SoftwarePkgReview
	Logs        []SoftwarePkgOperationLog

	ApprovedBy []SoftwarePkgApprover
	RejectedBy []SoftwarePkgApprover
}

func (entity *SoftwarePkgBasicInfo) Sig() string {
	return entity.Application.ImportingPkgSig.ImportingPkgSig()
}

func (entity *SoftwarePkgBasicInfo) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing()
}

func (entity *SoftwarePkgBasicInfo) AddReview(ur *UserReview) (bool, error) {
	if !entity.Phase.IsReviewing() || !entity.CI.isSuccess() {
		return false, errors.New("can't do this")
	}

	entity.Review.add(ur)

	b := entity.Review.pass()
	if b {
		entity.Phase = dp.PackagePhaseCreatingRepo
	}

	return b, nil
}

// TODO delete
func (entity *SoftwarePkgBasicInfo) ApproveBy(user *SoftwarePkgApprover) (bool, error) {
	return false, nil
}

func (entity *SoftwarePkgBasicInfo) RejectBy(user *Reviewer) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !user.isTC() {
		return allerror.NewNoPermission("not tc")
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.User, dp.PackageOperationLogActionReject, entity.Id,
		),
	)

	return nil
}

func (entity *SoftwarePkgBasicInfo) Abandon(user *User) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return notImporter
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionAbandon, entity.Id,
		),
	)

	return nil
}

func (entity *SoftwarePkgBasicInfo) RerunCI(user *User) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, incorrectPhase
	}

	if entity.CI.Status.IsCIRunning() {
		return false, allerror.New(allerror.ErrorCodeCIIsRunning, "ci is running")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return false, notImporter
	}

	if entity.CI.Status.IsCIWaiting() {
		return false, nil
	}

	entity.CI = SoftwarePkgCI{Status: dp.PackageCIStatusWaiting}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionResunci, entity.Id,
		),
	)

	return true, nil
}

func (entity *SoftwarePkgBasicInfo) UpdateApplication(cmd *SoftwarePkgApplication, user *User) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return notImporter
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

	Logs     []SoftwarePkgOperationLog
	Comments []SoftwarePkgReviewComment
}

func NewSoftwarePkg(user *User, name dp.PackageName, app *SoftwarePkgApplication) SoftwarePkgBasicInfo {
	return SoftwarePkgBasicInfo{
		PkgName:     name,
		Importer:    *user,
		Phase:       dp.PackagePhaseReviewing,
		CI:          SoftwarePkgCI{Status: dp.PackageCIStatusWaiting},
		Application: *app,
		AppliedAt:   utils.Now(),
	}
}
