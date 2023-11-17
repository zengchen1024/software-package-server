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
	ciInstance pkgCI

	notfound       = allerror.NewNotFound(allerror.ErrorCodePkgNotFound, "not found")
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

func (basic *SoftwarePkgBasicInfo) update(info *SoftwarePkgBasicInfo) []dp.PkgModificationCategory {
	categories := []dp.PkgModificationCategory{}

	if basic.Name.PackageName() != info.Name.PackageName() {
		categories = append(categories, dp.PkgModificationCategoryPkgName)
	}

	if basic.Desc.PackageDesc() != info.Desc.PackageDesc() {
		categories = append(categories, dp.PkgModificationCategoryPkgDesc)
	}

	if basic.Reason.ReasonToImportPkg() != info.Reason.ReasonToImportPkg() {
		categories = append(categories, dp.PkgModificationCategoryPkgReason)
	}

	if basic.Upstream.URL() != info.Upstream.URL() {
		categories = append(categories, dp.PkgModificationCategoryUpstream)
	}

	if len(categories) > 0 {
		*basic = *info
	}

	return categories
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
	Initialized bool
}

func (entity *SoftwarePkg) IsCommitter(user dp.Account) bool {
	return entity.Repo.isCommitter(user)
}

func (entity *SoftwarePkg) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing()
}

func (entity *SoftwarePkg) FilesToDownload() []SoftwarePkgCodeFile {
	return entity.Code.filesToDownload()
}

func (entity *SoftwarePkg) SaveDownloadedFiles(files []SoftwarePkgCodeFile) bool {
	return entity.Code.saveDownloadedFiles(files)
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
		categories = append(categories, dp.PkgModificationCategorySig)

		entity.Sig = sig
	}

	if v := entity.Repo.update(repo); v != nil {
		categories = append(categories, v)
	}

	if spec != nil || srpm != nil {
		entity.Code.update(spec, srpm)
		categories = append(categories, dp.PkgModificationCategoryCode)
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

func (entity *SoftwarePkg) StartInitialization(pr dp.URL) error {
	if !entity.Phase.IsCreatingRepo() || entity.Initialized {
		return errors.New("can't do this")
	}

	entity.CommunityPR = pr

	return nil
}

func (entity *SoftwarePkg) HandleInitialized(pr dp.URL) error {
	b := !entity.Phase.IsCreatingRepo() ||
		entity.CommunityPR == nil ||
		entity.CommunityPR.URL() != pr.URL()
	if b {
		return errors.New("can't do this")
	}

	entity.Initialized = true

	return nil
}

func (entity *SoftwarePkg) HandleAlreadyExisted() error {
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
		Sig:       sig,
		Repo:      *repo,
		Basic:     *basic,
		Importer:  importer.Account,
		CI:        SoftwarePkgCI{status: dp.PackageCIStatusWaiting},
		Phase:     dp.PackagePhaseReviewing,
		AppliedAt: utils.Now(),
	}

	pkg.Code.update(spec, srpm)

	return pkg
}
