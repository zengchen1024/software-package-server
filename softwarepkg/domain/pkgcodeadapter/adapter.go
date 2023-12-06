package pkgcodeadapter

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type PkgCodeAdapter interface {
	Download([]domain.SoftwarePkgCodeSourceFile, dp.PackageName) (bool, error)
	Clear(int, dp.PackageName) error
}
