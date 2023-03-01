package repositoryimpl

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

type SoftwarePkgDO struct {
	UUID            uuid.UUID      `gorm:"column:uuid;type:uuid"`
	Phase           string         `gorm:"column:phase"`
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
		ImportUser:      "", // TODO pkg.Importer.Account() is nil
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
