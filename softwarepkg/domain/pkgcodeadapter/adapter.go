package pkgcodeadapter

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type PkgCodeAdapter interface {
	DownloadCodes([]domain.SoftwarePkgCodeSourceFile, dp.PackageName) (bool, error)
	ClearCodes(dp.PackageName) error
	ClearAll(int, dp.PackageName) error
}
