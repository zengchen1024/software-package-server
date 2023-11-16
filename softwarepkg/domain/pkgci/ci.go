package pkgci

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type PkgCI interface {
	DownloadPkgCode(pkg *domain.SoftwarePkg) ([]domain.SoftwarePkgCodeInfo, error)
}
