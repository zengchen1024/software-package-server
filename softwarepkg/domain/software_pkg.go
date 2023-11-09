package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

const (
	gitee  = "gitee"
	github = "github"
)

var (
	notImporter    = allerror.NewNoPermission("not the importer")
	incorrectPhase = allerror.New(allerror.ErrorCodeIncorrectPhase, "incorrect phase")
)

type User struct {
	Email   dp.Email
	Account dp.Account

	GiteeID  string
	GithubID string
}

func (u *User) ApplyTo(p dp.PackagePlatform) bool {
	v := p.PackagePlatform()

	if v == gitee && u.GiteeID != "" {
		return true
	}

	if v == github && u.GithubID != "" {
		return true
	}

	return false
}

type SoftwarePkgBasicInfo struct {
	Name     dp.PackageName
	Desc     dp.PackageDesc
	Reason   dp.ReasonToImportPkg
	Upstream dp.URL
}

func (basic *SoftwarePkgBasicInfo) update(info *SoftwarePkgBasicInfo) []dp.CheckItemCategory {
	categories := []dp.CheckItemCategory{}

	if basic.Name.PackageName() != info.Name.PackageName() {
		categories = append(categories, dp.CheckItemCategoryPkgName)
	}

	if basic.Desc.PackageDesc() != info.Desc.PackageDesc() {
		categories = append(categories, dp.CheckItemCategoryPkgDesc)
	}

	if basic.Reason.ReasonToImportPkg() != info.Reason.ReasonToImportPkg() {
		categories = append(categories, dp.CheckItemCategoryPkgReason)
	}

	if basic.Upstream.URL() != info.Upstream.URL() {
		categories = append(categories, dp.CheckItemCategoryUpstream)
	}

	if len(categories) > 0 {
		*basic = *info
	}

	return categories
}

// SoftwarePkgCode
type SoftwarePkgCode struct {
	Spec SoftwarePkgCodeFile
	SRPM SoftwarePkgCodeFile
}

// SoftwarePkgCodeFile
type SoftwarePkgCodeFile struct {
	Src   dp.URL // Src is the url user inputed
	Local dp.URL // Local is the url that is the local address of the file
}

func (f *SoftwarePkgCodeFile) Name() string {
	// TODO
	return ""
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Platform   dp.PackagePlatform
	Committers []dp.Account
}

func (r *SoftwarePkgRepo) isCommitter(u dp.Account) bool {
	for i := range r.Committers {
		if dp.IsSameAccount(r.Committers[i], u) {
			return true
		}
	}

	return false
}

func (r *SoftwarePkgRepo) update(r1 *SoftwarePkgRepo) []dp.CheckItemCategory {
	b := false
	categories := []dp.CheckItemCategory{}

	if r.Platform.PackagePlatform() != r1.Platform.PackagePlatform() {
		b = true
	}

	if !r.isSameCommitters(r1) {
		categories = append(categories, dp.CheckItemCategoryCommitter)

		b = true
	}

	if b {
		*r = *r1
	}

	return categories
}

func (r *SoftwarePkgRepo) isSameCommitters(r1 *SoftwarePkgRepo) bool {
	if len(r.Committers) != len(r1.Committers) {
		return false
	}

	m := r.committersMap()

	for i := range r1.Committers {
		if !m[r1.Committers[i].Account()] {
			return false
		}
	}

	return true
}

func (r *SoftwarePkgRepo) committersMap() map[string]bool {
	m := map[string]bool{}

	for i := range r.Committers {
		m[r.Committers[i].Account()] = true
	}

	return m
}

// SoftwarePkgCI
type SoftwarePkgCI struct {
	Id        int // The Id of running CI
	status    dp.PackageCIStatus
	StartTime int64 // deal with the case that the ci is timeout
}

func (ci *SoftwarePkgCI) start() (int, error) {
	if !ci.status.IsCIWaiting() {
		return 0, errors.New("can't do this")
	}

	ci.status = dp.PackageCIStatusRunning
	ci.StartTime = utils.Now()
	ci.Id++

	return ci.Id, nil
}

func (ci *SoftwarePkgCI) rerun() error {
	s := ci.Status()

	if s.IsCIRunning() {
		return allerror.New(allerror.ErrorCodeCIIsRunning, "ci is running")
	}

	if s.IsCIWaiting() {
		return errors.New("duplicate operation")
	}

	ci.status = dp.PackageCIStatusWaiting
	ci.StartTime = utils.Now()

	return nil
}

func (ci *SoftwarePkgCI) done(ciId int, success bool) error {
	if ci.Id != ciId || !ci.Status().IsCIRunning() {
		return errors.New("ignore the ci result")
	}

	if success {
		ci.status = dp.PackageCIStatusPassed
	} else {
		ci.status = dp.PackageCIStatusFailed
	}

	return nil
}

func (ci *SoftwarePkgCI) isSuccess() bool {
	return ci.status != nil && ci.status.IsCIPassed()
}

func (ci *SoftwarePkgCI) Status() dp.PackageCIStatus {
	if ci.status == nil {
		return nil
	}

	s := ci.status

	// TODO config
	if s.IsCIWaiting() {
		if ci.StartTime+30 < utils.Now() {
			return dp.PackageCIStatusTimeout
		}

		return s
	}

	// TODO config
	if s.IsCIRunning() && ci.StartTime+300 < utils.Now() {
		return dp.PackageCIStatusTimeout
	}

	return s
}

// SoftwarePkg
type SoftwarePkg struct {
	Id       string
	Sig      dp.ImportingPkgSig
	Repo     SoftwarePkgRepo
	Code     SoftwarePkgCode
	Basic    SoftwarePkgBasicInfo
	Importer dp.Account

	CI          SoftwarePkgCI
	Logs        []SoftwarePkgOperationLog
	Phase       dp.PackagePhase
	Reviews     []UserReview
	AppliedAt   int64
	CommunityPR dp.URL
}

func (entity *SoftwarePkg) IsCommitter(user dp.Account) bool {
	return entity.Repo.isCommitter(user)
}

func (entity *SoftwarePkg) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing()
}

func (entity *SoftwarePkg) AddReview(ur *UserReview, items []CheckItem) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, incorrectPhase
	}

	if !entity.CI.isSuccess() {
		return false, allerror.New(
			allerror.ErrorCodeCIIsNotReady, "ci is not successful yet",
		)
	}

	if err := entity.addReview(ur, items); err != nil {
		return false, err
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			ur.Reviewer.Account, dp.PackageOperationLogActionReview,
		),
	)

	b := entity.doesPassReview(items)
	if b {
		entity.Phase = dp.PackagePhaseCreatingRepo
	}

	return b, nil
}

func (entity *SoftwarePkg) RejectBy(user *Reviewer) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	roles := maintainerInstance.Roles(entity, user)

	isTC := false

	for i := range roles {
		if roles[i].IsTC() {
			isTC = true
			break
		}
	}

	if !isTC {
		return allerror.NewNoPermission("not the tc")
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionReject,
		),
	)

	return nil
}

func (entity *SoftwarePkg) Abandon(user *User) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notImporter
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionAbandon,
		),
	)

	return nil
}

func (entity *SoftwarePkg) RerunCI(user *User) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notImporter
	}

	if err := entity.CI.rerun(); err != nil {
		return err
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionResunci,
		),
	)

	return nil
}

func (entity *SoftwarePkg) UpdateApplication(
	user *User,
	sig dp.ImportingPkgSig,
	repo *SoftwarePkgRepo,
	basic *SoftwarePkgBasicInfo,
	items []CheckItem,
) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notImporter
	}

	categories := entity.Basic.update(basic)

	if entity.Sig.ImportingPkgSig() != sig.ImportingPkgSig() {
		categories = append(categories, dp.CheckItemCategorySig)

		entity.Sig = sig
	}

	if v := entity.Repo.update(repo); len(v) > 0 {
		categories = append(categories, v...)
	}

	if len(categories) > 0 {
		entity.clearReview(categories, items)
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionUpdate,
		),
	)

	return nil
}

func (entity *SoftwarePkg) HandleCIChecking() (int, error) {
	if !entity.Phase.IsReviewing() {
		return 0, errors.New("can't do this")
	}

	return entity.CI.start()
}

func (entity *SoftwarePkg) HandleCIDone(ciId int, success bool) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	return entity.CI.done(ciId, success)
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
		Importer: importer.Account,

		CI:        SoftwarePkgCI{status: dp.PackageCIStatusWaiting},
		Phase:     dp.PackagePhaseReviewing,
		AppliedAt: utils.Now(),
	}
}
