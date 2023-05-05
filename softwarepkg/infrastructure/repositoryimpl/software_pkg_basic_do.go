package repositoryimpl

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/plugin/optimisticlock"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

const (
	fieldId              = "uuid"
	fieldPhase           = "phase"
	fieldVersion         = "version"
	fieldImporter        = "importer"
	fieldAppliedAt       = "applied_at"
	fieldPackageName     = "package_name"
	fieldPackagePlatform = "package_platform"

	frozenStatus   = "frozen"
	unfrozenStatus = "unfrozen"
)

func (s softwarePkgBasic) toSoftwarePkgBasicDO(pkg *domain.SoftwarePkgBasicInfo, do *SoftwarePkgBasicDO) (err error) {
	email, err := toEmailDO(pkg.Importer.Email)
	if err != nil {
		return err
	}

	app := &pkg.Application

	*do = SoftwarePkgBasicDO{
		Id:              uuid.New(),
		PackageName:     pkg.PkgName.PackageName(),
		Importer:        pkg.Importer.Account.Account(),
		ImporterEmail:   email,
		Phase:           pkg.Phase.PackagePhase(),
		CIStatus:        pkg.CI.Status.PackageCIStatus(),
		SpecURL:         app.SourceCode.SpecURL.URL(),
		SrcRPMURL:       app.SourceCode.SrcRPMURL.URL(),
		PackageDesc:     app.PackageDesc.PackageDesc(),
		PackagePlatform: app.PackagePlatform.PackagePlatform(),
		Sig:             app.ImportingPkgSig.ImportingPkgSig(),
		ReasonToImport:  app.ReasonToImportPkg.ReasonToImportPkg(),
		AppliedAt:       pkg.AppliedAt,
		UpdatedAt:       pkg.AppliedAt,
		ApprovedBy:      toStringArray(pkg.ApprovedBy),
		RejectedBy:      toStringArray(pkg.RejectedBy),
	}

	if pkg.Frozen {
		do.Frozen = frozenStatus
	} else {
		do.Frozen = unfrozenStatus
	}

	if pkg.RepoLink != nil {
		do.RepoLink = pkg.RepoLink.URL()
	}

	if pkg.RelevantPR != nil {
		do.RelevantPR = pkg.RelevantPR.URL()
	}

	return
}

type SoftwarePkgBasicDO struct {
	// must set "uuid" as the name of column
	Id              uuid.UUID              `gorm:"column:uuid;type:uuid"`
	PackageName     string                 `gorm:"column:package_name"`
	Importer        string                 `gorm:"column:importer"`
	ImporterEmail   string                 `gorm:"column:importer_email"`
	RepoLink        string                 `gorm:"column:repo_link"`
	Phase           string                 `gorm:"column:phase"`
	CIStatus        string                 `gorm:"column:ci_status"`
	SpecURL         string                 `gorm:"column:spec_url"`
	SrcRPMURL       string                 `gorm:"column:src_rpm_url"`
	PackageDesc     string                 `gorm:"column:package_desc"`
	PackagePlatform string                 `gorm:"column:package_platform"`
	RelevantPR      string                 `gorm:"column:relevant_pr"`
	Sig             string                 `gorm:"column:sig"`
	Frozen          string                 `gorm:"column:frozen"`
	ReasonToImport  string                 `gorm:"column:reason_to_import"`
	ApprovedBy      pq.StringArray         `gorm:"column:approvedby;type:text[];default:'{}'"`
	RejectedBy      pq.StringArray         `gorm:"column:rejectedby;type:text[];default:'{}'"`
	AppliedAt       int64                  `gorm:"column:applied_at"`
	UpdatedAt       int64                  `gorm:"column:updated_at"`
	Version         optimisticlock.Version `gorm:"column:version"`
}

func (do *SoftwarePkgBasicDO) toSoftwarePkgBasicInfo() (info domain.SoftwarePkgBasicInfo, err error) {
	info.Id = do.Id.String()

	if info.PkgName, err = dp.NewPackageName(do.PackageName); err != nil {
		return
	}

	if do.RepoLink != "" {
		if info.RepoLink, err = dp.NewURL(do.RepoLink); err != nil {
			return
		}
	}

	if do.RelevantPR != "" {
		if info.RelevantPR, err = dp.NewURL(do.RelevantPR); err != nil {
			return
		}
	}

	if info.Importer.Account, err = dp.NewAccount(do.Importer); err != nil {
		return
	}

	if info.Importer.Email, err = toEmail(do.ImporterEmail); err != nil {
		return
	}

	if info.Phase, err = dp.NewPackagePhase(do.Phase); err != nil {
		return
	}

	info.AppliedAt = do.AppliedAt

	if err = do.toSoftwarePkgApplication(&info.Application); err != nil {
		return
	}

	if info.ApprovedBy, err = do.toAccounts(do.ApprovedBy); err != nil {
		return
	}

	if info.CI.Status, err = dp.NewPackageCIStatus(do.CIStatus); err != nil {
		return
	}

	info.RejectedBy, err = do.toAccounts(do.RejectedBy)

	info.Frozen = do.Frozen == frozenStatus

	return
}

func (do *SoftwarePkgBasicDO) toAccounts(v []string) (r []dp.Account, err error) {
	if len(v) == 0 {
		return
	}

	r = make([]dp.Account, len(v))
	for i := range v {
		if r[i], err = dp.NewAccount(v[i]); err != nil {
			return
		}
	}

	return
}

func (do *SoftwarePkgBasicDO) toSoftwarePkgApplication(app *domain.SoftwarePkgApplication) (err error) {
	if app.ReasonToImportPkg, err = dp.NewReasonToImportPkg(do.ReasonToImport); err != nil {
		return
	}

	if app.PackageDesc, err = dp.NewPackageDesc(do.PackageDesc); err != nil {
		return
	}

	if app.PackagePlatform, err = dp.NewPackagePlatform(do.PackagePlatform); err != nil {
		return
	}

	if app.ImportingPkgSig, err = dp.NewImportingPkgSig(do.Sig); err != nil {
		return
	}

	if app.SourceCode.SrcRPMURL, err = dp.NewURL(do.SrcRPMURL); err != nil {
		return
	}

	app.SourceCode.SpecURL, err = dp.NewURL(do.SpecURL)

	return
}

func toStringArray(v []dp.Account) (arr pq.StringArray) {
	arr = make(pq.StringArray, len(v))
	for i, account := range v {
		arr[i] = account.Account()
	}

	return
}

func toEmailDO(e dp.Email) (string, error) {
	return utils.Encryption.Encrypt([]byte(e.Email()))
}

func toEmail(e string) (dp.Email, error) {
	v, err := utils.Encryption.Decrypt(e)
	if err != nil {
		return nil, err
	}

	return dp.NewEmail(string(v))
}
