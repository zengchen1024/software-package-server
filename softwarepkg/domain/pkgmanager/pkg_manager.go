package pkgmanager

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

type PkgManager interface {
	IsPkgExisted(dp.PackageName) bool
}
