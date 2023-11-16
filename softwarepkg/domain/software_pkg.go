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
	ciInstance     pkgCI
	notfound       = allerror.NewNotFound(allerror.ErrorCodePkgNotFound, "not found")
	incorrectPhase = allerror.New(allerror.ErrorCodeIncorrectPhase, "incorrect phase")
)

type pkgCI interface {
	// return new pr num
	StartNewCI(pkg *SoftwarePkg) (int, error)
	ClearCI(pkg *SoftwarePkg) error
}

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

func (code *SoftwarePkgCode) isReady() bool {
	return code.Spec.isReady() && code.SRPM.isReady()
}

func (code *SoftwarePkgCode) update(spec, srpm dp.URL) {
	if spec != nil {
		code.Spec.update(spec)
	}

	if srpm != nil {
		code.SRPM.update(srpm)
	}
}

func (code *SoftwarePkgCode) filesToDownload() []SoftwarePkgCodeInfo {
	r := []SoftwarePkgCodeInfo{}

	if !code.Spec.isReady() {
		r = append(r, code.Spec.SoftwarePkgCodeInfo)
	}

	if !code.SRPM.isReady() {
		r = append(r, code.SRPM.SoftwarePkgCodeInfo)
	}

	return r
}

func (code *SoftwarePkgCode) saveDownloadedFiles(infos []SoftwarePkgCodeInfo) (changed bool) {
	spec := false

	for i := range infos {
		item := &infos[i]

		if !spec && code.Spec.saveLocalPath(item) {
			changed = true
		}

		if code.SRPM.saveLocalPath(item) {
			changed = true
		}
	}

	return
}

// SoftwarePkgCodeInfo
type SoftwarePkgCodeInfo struct {
	Src       dp.URL // Src is the url user inputed
	Local     dp.URL // Local is the url that is the local address of the file
	UpdatedAt int64
}

func (f *SoftwarePkgCodeInfo) FileName() string {
	return f.Src.FileName()
}

func (info *SoftwarePkgCodeInfo) isSame(info1 *SoftwarePkgCodeInfo) bool {
	return info.UpdatedAt == info1.UpdatedAt && info.Src.URL() == info1.Src.URL()
}

// SoftwarePkgCodeFile
type SoftwarePkgCodeFile struct {
	SoftwarePkgCodeInfo

	Dirty bool // if true, the code should be updated.
	//Reason string // the reason why can't download the code file
}

func (f *SoftwarePkgCodeFile) isReady() bool {
	return f.Local != nil && !f.Dirty
}

func (f *SoftwarePkgCodeFile) update(src dp.URL) {
	f.Src = src
	f.Local = nil
	f.Dirty = true
	// f.Reason = ""
	f.UpdatedAt = utils.Now()
}

func (f *SoftwarePkgCodeFile) saveLocalPath(info *SoftwarePkgCodeInfo) bool {
	if !f.isReady() && f.isSame(info) {
		f.Local = info.Local
		f.Dirty = false

		return true
	}

	return false
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

func (r *SoftwarePkgRepo) update(r1 *SoftwarePkgRepo) dp.CheckItemCategory {
	if p := r1.Platform; p != nil && r.Platform.PackagePlatform() != p.PackagePlatform() {
		r.Platform = p
	}

	if len(r1.Committers) > 0 && !r.isSameCommitters(r1) {
		r.Committers = r1.Committers

		return dp.CheckItemCategoryCommitter
	}

	return nil
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

func (ci *SoftwarePkgCI) start(pkg *SoftwarePkg) error {
	if !ci.status.IsCIWaiting() {
		return errors.New("can't do this")
	}

	ci.status = dp.PackageCIStatusRunning
	ci.StartTime = utils.Now()

	v, err := ciInstance.StartNewCI(pkg)
	if err == nil {
		ci.Id = v
	}

	return err
}

func (ci *SoftwarePkgCI) retest() error {
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

	ci.StartTime = 0

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

func (entity *SoftwarePkg) FilesToDownload() []SoftwarePkgCodeInfo {
	return entity.Code.filesToDownload()
}

func (entity *SoftwarePkg) SaveDownloadedFiles(infos []SoftwarePkgCodeInfo) bool {
	return entity.Code.saveDownloadedFiles(infos)
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
	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notfound
	}

	if !entity.Phase.IsReviewing() {
		return incorrectPhase
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

func (entity *SoftwarePkg) Retest(user *User) error {
	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notfound
	}

	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !entity.Code.isReady() {
		return allerror.New(allerror.ErrorCodePkgCodeNotReady, "code not ready")
	}

	if err := entity.CI.retest(); err != nil {
		return err
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionRetest,
		),
	)

	return nil
}

func (entity *SoftwarePkg) UpdateApplication(
	user *User,
	sig dp.ImportingPkgSig,
	repo *SoftwarePkgRepo,
	basic *SoftwarePkgBasicInfo,
	spec dp.URL,
	srpm dp.URL,
) error {
	if !dp.IsSameAccount(user.Account, entity.Importer) {
		return notfound
	}

	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	categories := entity.Basic.update(basic)

	if entity.Sig.ImportingPkgSig() != sig.ImportingPkgSig() {
		categories = append(categories, dp.CheckItemCategorySig)

		entity.Sig = sig
	}

	if v := entity.Repo.update(repo); v != nil {
		categories = append(categories, v)
	}

	if spec != nil || srpm != nil {
		entity.Code.update(spec, srpm)
		categories = append(categories, dp.CheckItemCategoryCode)
	}

	if len(categories) == 0 {
		return errors.New("nothing changed")
	}

	// TODO
	var items []CheckItem

	entity.clearReview(categories, items)

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionUpdate,
		),
	)

	return nil
}

func (entity *SoftwarePkg) StartCI() error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	return entity.CI.start(entity)
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
		return incorrectPhase
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

func (entity *SoftwarePkg) HandleRepoCodePushed() error {
	if !entity.Phase.IsCreatingRepo() {
		return incorrectPhase
	}

	entity.Phase = dp.PackagePhaseImported

	return nil
}

func NewSoftwarePkg(
	sig dp.ImportingPkgSig,
	repo *SoftwarePkgRepo,
	basic *SoftwarePkgBasicInfo,
	spec dp.URL,
	srpm dp.URL,
	importer *User,
) SoftwarePkg {
	pkg := SoftwarePkg{
		Sig:      sig,
		Repo:     *repo,
		Basic:    *basic,
		Importer: importer.Account,

		CI: SoftwarePkgCI{
			status: dp.PackageCIStatusWaiting,
		},
		Phase:     dp.PackagePhaseReviewing,
		AppliedAt: utils.Now(),
	}

	pkg.Code.update(spec, srpm)

	return pkg
}
