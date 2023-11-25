package domain

import (
	"errors"
	"fmt"

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
	Importer dp.Account

	CI          SoftwarePkgCI
	Logs        []SoftwarePkgOperationLog
	Phase       dp.PackagePhase
	Reviews     []UserReview
	AppliedAt   int64
	CommunityPR dp.URL
	Initialized bool
}

func (entity *SoftwarePkg) CheckItems() []CheckItem {
	other := entity.otherCheckItems()

	r := make([]CheckItem, 0, len(other)+len(commonCheckItems))
	r = append(r, commonCheckItems...) // don't change commonCheckItems by copy it.
	r = append(r, other...)

	return r
}

func (entity *SoftwarePkg) otherCheckItems() []CheckItem {
	v := []CheckItem{
		{
			Id:            entity.Sig.ImportingPkgSig(),
			Name:          "Sig",
			Desc:          fmt.Sprintf("软件包被%s Sig接纳", entity.Sig.ImportingPkgSig()),
			Owner:         dp.CommunityRoleSigMaintainer,
			OnlyOwner:     true,
			Modifications: []string{pkgModificationSig},
		},
	}

	for i := range entity.Repo.Committers {
		c := entity.Repo.Committers[i].Account.Account()

		if c == entity.Importer.Account() {
			continue
		}

		v = append(v, CheckItem{
			Id:            c,
			Name:          "软件包维护人",
			Desc:          fmt.Sprintf("%s 同意作为此软件包的维护人", c),
			Owner:         dp.CommunityRoleCommitter,
			Keep:          true,
			OnlyOwner:     true,
			Modifications: []string{pkgModificationCommitter},
		})
	}

	return v
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

func (entity *SoftwarePkg) RepoLink() string {
	return entity.Repo.repoLink(entity.Basic.Name)
}

func (entity *SoftwarePkg) FilesToDownload() []SoftwarePkgCodeSourceFile {
	return entity.Code.filesToDownload()
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

func (entity *SoftwarePkg) RejectBy(user *Reviewer) error {
	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	if tc, _ := maintainerInstance.Roles(entity, user); !tc {
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

func (entity *SoftwarePkg) Abandon(user dp.Account) error {
	if !dp.IsSameAccount(user, entity.Importer) {
		return notfound
	}

	if !entity.Phase.IsReviewing() {
		return incorrectPhase
	}

	entity.Phase = dp.PackagePhaseClosed

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user, dp.PackageOperationLogActionAbandon,
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

func (entity *SoftwarePkg) Update(user *User, info *SoftwarePkgUpdateInfo) error {
	if !dp.IsSameAccount(user.Account, entity.Importer) {
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
	}

	if info.Spec != nil || info.SRPM != nil {
		entity.Code.update(info.Spec, info.SRPM)
		ms = append(ms, pkgModificationCode)
	}

	if len(ms) == 0 {
		return errors.New("nothing changed")
	}

	items := append(entity.otherCheckItems(), commonCheckItems...)
	entity.clearReview(ms, items)

	entity.Logs = append(
		entity.Logs,
		NewSoftwarePkgOperationLog(
			user.Account, dp.PackageOperationLogActionUpdate,
		),
	)

	/*
		if entity.doesPassReview(items) {
			entity.Phase = dp.PackagePhaseCreatingRepo
		}
	*/

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
		CI:        SoftwarePkgCI{status: dp.PackageCIStatusWaiting},
		Sig:       sig,
		Repo:      *repo,
		Basic:     *basic,
		Phase:     dp.PackagePhaseReviewing,
		Importer:  importer.Account,
		AppliedAt: utils.Now(),
	}

	pkg.Code.update(spec, srpm)

	return pkg
}
