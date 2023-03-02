package repositoryimpl

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	applyTime = "apply_time"
)

type SoftwarePkgDO struct {
	UUID            uuid.UUID      `gorm:"column:uuid;type:uuid"`
	Phase           string         `gorm:"column:phase"`
	Review          string         `gorm:"column:review"`
	SourceCode      string         `gorm:"column:source_code"`
	PackageSig      string         `gorm:"column:package_sig"`
	ImportUser      string         `gorm:"column:import_user"`
	PackageName     string         `gorm:"column:package_name"`
	PackageDesc     string         `gorm:"column:package_desc"`
	PackageReason   string         `gorm:"column:package_reason"`
	PackageLicense  string         `gorm:"column:package_license"`
	PackagePlatform string         `gorm:"column:package_platform"`
	PackageRepoLink string         `gorm:"column:package_repo_link"`
	RejectUser      pq.StringArray `gorm:"column:reject_user;type:text[];default:'{}'"`
	ApproveUser     pq.StringArray `gorm:"column:approve_user;type:text[];default:'{}'"`
	Version         int            `gorm:"column:version"`
	AppliedAt       int64          `gorm:"column:apply_time"`
	UpdatedAt       int64          `gorm:"column:update_time"`
}

func (s softwarePkgTable) toSoftwarePkgDO(pkg *domain.SoftwarePkgBasicInfo, do *SoftwarePkgDO) {
	*do = SoftwarePkgDO{
		UUID:            uuid.New(),
		Phase:           pkg.Phase.PackagePhase(),
		SourceCode:      pkg.Application.SourceCode.Address.URL(),
		ImportUser:      "ceshi", // TODO pkg.Importer.Account() is nil
		PackageSig:      pkg.Application.ImportingPkgSig.ImportingPkgSig(),
		PackageName:     pkg.Application.PackageName.PackageName(),
		PackageDesc:     pkg.Application.PackageDesc.PackageDesc(),
		PackageReason:   pkg.Application.ReasonToImportPkg.ReasonToImportPkg(),
		PackageLicense:  pkg.Application.SourceCode.License.License(),
		PackagePlatform: pkg.Application.PackagePlatform.PackagePlatform(),
		AppliedAt:       pkg.AppliedAt,
		UpdatedAt:       pkg.AppliedAt,
	}
}

func (s SoftwarePkgDO) toSoftwarePkgBasicInfo() (info domain.SoftwarePkgBasicInfo, err error) {
	info.Id = s.UUID.String()

	if info.PkgName, err = dp.NewPackageName(s.PackageName); err != nil {
		return
	}

	if s.PackageRepoLink != "" {
		if info.RepoLink, err = dp.NewURL(s.PackageRepoLink); err != nil {
			return
		}
	}

	if info.Importer, err = dp.NewAccount(s.ImportUser); err != nil {
		return
	}

	if info.Phase, err = dp.NewPackagePhase(s.Phase); err != nil {
		return
	}

	info.AppliedAt = s.AppliedAt

	if err = s.toSoftwarePkgApplication(&info.Application); err != nil {
		return
	}

	if info.ApprovedBy, err = s.toAccounts(s.ApproveUser); err != nil {
		return
	}

	info.RejectedBy, err = s.toAccounts(s.RejectUser)

	return
}

func (s SoftwarePkgDO) toAccounts(v []string) (r []dp.Account, err error) {
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

func (s SoftwarePkgDO) toSoftwarePkgApplication(app *domain.SoftwarePkgApplication) (err error) {
	if app.ReasonToImportPkg, err = dp.NewReasonToImportPkg(s.PackageReason); err != nil {
		return
	}

	if app.PackageName, err = dp.NewPackageName(s.PackageName); err != nil {
		return
	}

	if app.PackageDesc, err = dp.NewPackageDesc(s.PackageDesc); err != nil {
		return
	}

	if app.PackagePlatform, err = dp.NewPackagePlatform(s.PackagePlatform); err != nil {
		return
	}

	if app.ImportingPkgSig, err = dp.NewImportingPkgSig(s.PackageSig); err != nil {
		return
	}

	if app.SourceCode.License, err = dp.NewLicense(s.PackageLicense); err != nil {
		return
	}

	app.SourceCode.Address, err = dp.NewURL(s.SourceCode)

	return
}
