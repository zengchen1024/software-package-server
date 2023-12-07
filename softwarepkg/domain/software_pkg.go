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
	incorrectPhase = allerror.New(allerror.ErrorCodePkgIncorrectPhase, "incorrect phase")
)

type SoftwarePkgUpdateInfo struct {
	Sig      dp.ImportingPkgSig
	Repo     SoftwarePkgRepo
	Spec     dp.URL
	SRPM     dp.URL
	Desc     dp.PackageDesc
	Purpose  dp.PurposeToImportPkg
	Upstream dp.URL
}

type User struct {
	Email   dp.Email
	Account dp.Account

	GiteeID  string
	GithubID string
}

func (u *User) Id(p string) string {
	switch p {
	case gitee:
		return u.GiteeID

	case github:
		return u.GithubID

	default:
		return ""
	}
}

type SoftwarePkgBasicInfo struct {
	Name     dp.PackageName
	Desc     dp.PackageDesc
	Purpose  dp.PurposeToImportPkg
	Upstream dp.URL
}

func (basic *SoftwarePkgBasicInfo) update(info *SoftwarePkgUpdateInfo) []string {
	ms := []string{}

	if v := info.Desc; v != nil && v.PackageDesc() != basic.Desc.PackageDesc() {
		basic.Desc = v

		ms = append(ms, pkgModificationPkgDesc)
	}

	if v := info.Purpose; v != nil && v.PurposeToImportPkg() != basic.Purpose.PurposeToImportPkg() {
		basic.Purpose = v

		ms = append(ms, pkgModificationPurpose)
	}

	if v := info.Upstream; v != nil && v.URL() != basic.Upstream.URL() {
		basic.Upstream = v

		ms = append(ms, pkgModificationUpstream)
	}

	return ms
}

// SoftwarePkg
type SoftwarePkg struct {
	Id       string
	Sig      dp.ImportingPkgSig
	Repo     SoftwarePkgRepo
	Code     SoftwarePkgCode
	Basic    SoftwarePkgBasicInfo
	Importer PkgCommitter

	CI          SoftwarePkgCI
	Logs        []SoftwarePkgOperationLog
	Phase       dp.PackagePhase
	Reviews     []UserReview
	AppliedAt   int64
	CommunityPR dp.URL
	Initialized bool
}

func (entity *SoftwarePkg) isCommitter(user dp.Account) bool {
	return entity.Repo.isCommitter(user)
}

func (entity *SoftwarePkg) CanAddReviewComment() error {
	if entity.Phase.IsReviewing() {
		return nil
	}

	return incorrectPhase
}

func (entity *SoftwarePkg) PackageName() dp.PackageName {
	return entity.Basic.Name
}

func (entity *SoftwarePkg) CIId() int {
	return entity.CI.Id
}

func (entity *SoftwarePkg) RepoLink() string {
	return entity.Repo.repoLink(entity.Basic.Name)
}

func (entity *SoftwarePkg) FilesToDownload() []SoftwarePkgCodeSourceFile {
	if entity.Phase.IsReviewing() {
		return entity.Code.filesToDownload()
	}

	return nil
}

func (entity *SoftwarePkg) SaveDownloadedFiles(files []SoftwarePkgCodeSourceFile) bool {
	return entity.Code.saveDownloadedFiles(files)
}

func (entity *SoftwarePkg) AddReview(ur *UserReview) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if !entity.CI.isSuccess() {
		return allerror.New(
			allerror.ErrorCodeCIIsNotReady, "ci is not successful yet",
		)
	}

	items := append(entity.otherCheckItems(), commonCheckItems...)

	if err := entity.addReview(ur, items); err != nil {
		return err
	}

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			ur.Reviewer.Account, dp.PackageOperationLogActionReview,
		),
	)

	if entity.doesPassReview(items) {
		entity.Phase = dp.PackagePhaseCreatingRepo
	}

	return nil
}

func (entity *SoftwarePkg) Close(user *Reviewer) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	var action dp.PackageOperationLogAction

	if entity.Importer.isMe(user.Account) {
		action = dp.PackageOperationLogActionAbandon
	} else {
		if tc, _ := maintainerInstance.Roles(entity, user); !tc {
			return notfound
		}

		action = dp.PackageOperationLogActionReject
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs, NewSoftwarePkgOperationLog(user.Account, action),
	)

	return nil
}

func (entity *SoftwarePkg) Retest(user *User) error {
	if !entity.Importer.isMe(user.Account) {
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

func (entity *SoftwarePkg) Update(importer *PkgCommitter, info *SoftwarePkgUpdateInfo) error {
	if !entity.Importer.isMe(importer.Account) {
		return notfound
	}

	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	ms := entity.Basic.update(info)

	if v := info.Sig; v.ImportingPkgSig() != entity.Sig.ImportingPkgSig() {
		entity.Sig = v

		ms = append(ms, pkgModificationSig)
	}

	if v := entity.Repo.update(&info.Repo); v != "" {
		ms = append(ms, v)

		entity.Importer.PlatformId = importer.PlatformId
	}

	if info.Spec != nil || info.SRPM != nil {
		entity.Code.update(info.Spec, info.SRPM)
		ms = append(ms, pkgModificationCode)
	}

	if len(ms) == 0 {
		return allerror.New(allerror.ErrorCodePkgNothingChanged, "nothing changed")
	}

	items := append(entity.otherCheckItems(), commonCheckItems...)
	entity.clearReview(ms, items)

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			importer.Account, dp.PackageOperationLogActionUpdate,
		),
	)

	if entity.doesPassReview(items) {
		entity.Phase = dp.PackagePhaseCreatingRepo
	}

	return nil
}

func (entity *SoftwarePkg) StartCI() error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	return entity.CI.start(entity)
}

func (entity *SoftwarePkg) HandleCIDone(ciId int, success bool) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
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
	importer *PkgCommitter,
) SoftwarePkg {
	now := utils.Now()

	pkg := SoftwarePkg{
		CI:        SoftwarePkgCI{status: dp.PackageCIStatusWaiting, StartTime: now},
		Sig:       sig,
		Repo:      *repo,
		Basic:     *basic,
		Phase:     dp.PackagePhaseReviewing,
		Importer:  *importer,
		AppliedAt: now,
	}

	pkg.Code.update(spec, srpm)

	return pkg
}
