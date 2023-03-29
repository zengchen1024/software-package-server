package pkgmanager

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type PkgManager interface {
	IsPkgExisted(dp.PackageName) bool
	GetPkg(dp.PackageName) (domain.SoftwarePkgBasicInfo, error)
}
