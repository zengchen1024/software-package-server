package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

type User struct {
	Email   dp.Email
	Account dp.Account

	GiteeID string
}

type SoftwarePkgBasicInfo struct {
	Name     dp.PackageName
	Desc     dp.PackageDesc
	Reason   dp.ReasonToImportPkg
	Upstream dp.URL
}

// SoftwarePkgCode
type SoftwarePkgCode struct {
	Spec SoftwarePkgCodeFile
	SRPM SoftwarePkgCodeFile
}

// SoftwarePkgCodeFile
type SoftwarePkgCodeFile struct {
	Src  dp.URL // Src is the url user inputed
	Path dp.URL // Path is the url that the actual storage address of the file
}

func (f *SoftwarePkgCodeFile) Name() string {
	// TODO
	return ""
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Link       dp.URL
	Platform   dp.PackagePlatform
	Committers []string
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

// SoftwarePkg
type SoftwarePkg struct {
	Id       string
	Sig      dp.ImportingPkgSig
	Repo     SoftwarePkgRepo
	Code     SoftwarePkgCode
	Basic    SoftwarePkgBasicInfo
	Importer User

	CI          SoftwarePkgCI
	Logs        []SoftwarePkgOperationLog
	Phase       dp.PackagePhase
	Review      SoftwarePkgReview
	AppliedAt   int64
	ApprovedBy  []SoftwarePkgApprover
	RejectedBy  []SoftwarePkgApprover
	CommunityPR dp.URL
}

func (entity *SoftwarePkg) ReviewResult() dp.PackageReviewResult {
	if len(entity.RejectedBy) > 0 {
		return dp.PkgReviewResultRejected
	}

	if entity.hasPassedReview() {
		return dp.PkgReviewResultApproved
	}

	return nil
}

func (entity *SoftwarePkg) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing()
}

func (entity *SoftwarePkg) ApproveBy(user *SoftwarePkgApprover) (bool, error) {
	if !entity.Phase.IsReviewing() || !entity.CI.isSuccess() {
		return false, errors.New("can't do this")
	}

	entity.ApprovedBy = append(entity.ApprovedBy, *user)

	b := entity.hasPassedReview()
	if b {
		entity.Phase = dp.PackagePhaseCreatingRepo
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionApprove, entity.Id,
		),
	)

	return b, nil
}

func (entity *SoftwarePkg) hasPassedReview() bool {
	if entity.Sig.ImportingPkgSig() == config.EcopkgSig {
		return len(entity.ApprovedBy) > 0
	}

	numApprovedByTc := 0
	numApprovedBySigMaintainer := 0

	for i := range entity.ApprovedBy {
		if entity.ApprovedBy[i].IsTC {
			numApprovedByTc++
			numApprovedBySigMaintainer++
		} else {
			numApprovedBySigMaintainer++
		}
	}

	c := numApprovedByTc >= config.MinNumApprovedByTC
	c1 := numApprovedBySigMaintainer >= config.MinNumApprovedBySigMaintainer

	return c && c1
}

func (entity *SoftwarePkg) RejectBy(user *SoftwarePkgApprover) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	entity.RejectedBy = append(entity.RejectedBy, *user)

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionReject, entity.Id,
		),
	)

	return nil
}

func (entity *SoftwarePkg) Abandon(user *User) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return errorNotTheImporter
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

func (entity *SoftwarePkg) RerunCI(user *User) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, errors.New("can't do this")
	}

	if entity.CI.Status.IsCIRunning() {
		return false, errorCIIsRunning
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return false, errorNotTheImporter
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

func (entity *SoftwarePkg) UpdateApplication(
	basic *SoftwarePkgBasicInfo,
	sig dp.ImportingPkgSig,
	repo *SoftwarePkgRepo,
	user *User,
) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user.Account, entity.Importer.Account) {
		return errorNotTheImporter
	}

	// TODO

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionUpdate, entity.Id,
		),
	)

	return nil
}

func (entity *SoftwarePkg) HandleCIChecking() error {
	b := entity.Phase.IsReviewing() && entity.CI.Status.IsCIWaiting()
	if !b {
		return errors.New("can't do this")
	}

	entity.CI.Status = dp.PackageCIStatusRunning
	//entity.CI.startTime = utils.Now()

	return nil
}

func (entity *SoftwarePkg) HandleCIChecked(success bool, prNum int) error {
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

func (entity *SoftwarePkg) HandlePkgInitialized(pr dp.URL) error {
	if !entity.Phase.IsCreatingRepo() {
		return errors.New("can't do this")
	}

	entity.CommunityPR = pr

	return nil
}

func (entity *SoftwarePkg) HandlePkgAlreadyExisted() error {
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

func (entity *SoftwarePkg) HandleRepoCreated(info RepoCreatedInfo) error {
	return nil
}

func (entity *SoftwarePkg) HandleCodeSaved(info RepoCreatedInfo) error {
	if err := entity.HandleRepoCreated(info); err != nil {
		return err
	}

	entity.Phase = dp.PackagePhaseImported

	return nil
}

func NewSoftwarePkg(
	sig dp.ImportingPkgSig,
	repo *SoftwarePkgRepo,
	code *SoftwarePkgCode,
	basic *SoftwarePkgBasicInfo,
	importer *User,
) SoftwarePkg {
	return SoftwarePkg{
		Sig:      sig,
		Repo:     *repo,
		Code:     *code,
		Basic:    *basic,
		Importer: *importer,

		CI:        SoftwarePkgCI{Status: dp.PackageCIStatusWaiting},
		Phase:     dp.PackagePhaseReviewing,
		AppliedAt: utils.Now(),
	}
}
