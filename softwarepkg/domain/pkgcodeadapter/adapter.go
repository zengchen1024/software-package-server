package pkgcodeadapter

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type PkgCodeAdapter interface {
	Download([]domain.SoftwarePkgCodeSourceFile, dp.PackageName) error
	Clear(pkg *domain.SoftwarePkg) error
}
