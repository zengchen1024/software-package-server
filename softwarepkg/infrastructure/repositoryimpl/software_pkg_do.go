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
	ApplyTime       int64          `gorm:"column:apply_time"`
	UpdateTime      int64          `gorm:"column:update_time"`
}

func (SoftwarePkgDO) TableName() string {
	return "software_pkg"
}

func (s softwarePkgImpl) toSoftwarePkgDO(pkg *domain.SoftwarePkgBasicInfo) *SoftwarePkgDO {
	softwareDO := &SoftwarePkgDO{
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
		ApplyTime:       pkg.AppliedAt,
		UpdateTime:      pkg.AppliedAt,
	}

	return softwareDO
}

func (s SoftwarePkgDO) toSoftwarePkgSummary() (info domain.SoftwarePkgBasicInfo, err error) {
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

	info.AppliedAt = s.ApplyTime

	var pkg domain.SoftwarePkgApplication
	if pkg, err = s.toSoftwarePkgApplication(); err != nil {
		return
	} else {
		info.Application = pkg
	}

	for _, v := range s.ApproveUser {
		var approve dp.Account
		if approve, err = dp.NewAccount(v); err != nil {
			return
		}
		info.ApprovedBy = append(info.ApprovedBy, approve)
	}

	for _, v := range s.RejectUser {
		var reject dp.Account
		if reject, err = dp.NewAccount(v); err != nil {
			return
		}
		info.RejectedBy = append(info.RejectedBy, reject)
	}

	return
}

func (s SoftwarePkgDO) toSoftwarePkgApplication() (pkg domain.SoftwarePkgApplication, err error) {
	if pkg.ReasonToImportPkg, err = dp.NewReasonToImportPkg(s.PackageReason); err != nil {
		return
	}

	if pkg.PackageName, err = dp.NewPackageName(s.PackageName); err != nil {
		return
	}

	if pkg.PackageDesc, err = dp.NewPackageDesc(s.PackageDesc); err != nil {
		return
	}

	if pkg.PackagePlatform, err = dp.NewPackagePlatform(s.PackagePlatform); err != nil {
		return
	}

	if pkg.ImportingPkgSig, err = dp.NewImportingPkgSig(s.PackageSig); err != nil {
		return
	}

	if pkg.SourceCode.License, err = dp.NewLicense(s.PackageLicense); err != nil {
		return
	}

	pkg.SourceCode.Address, err = dp.NewURL(s.SourceCode)

	return
}
