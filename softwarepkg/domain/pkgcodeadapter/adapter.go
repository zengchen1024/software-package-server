package pkgcodeadapter

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type PkgCodeAdapter interface {
	Download(pkg *domain.SoftwarePkg) ([]domain.SoftwarePkgCodeFile, error)
	Clear(pkg *domain.SoftwarePkg) error
}
