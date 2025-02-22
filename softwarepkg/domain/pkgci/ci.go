package pkgci

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type PkgCI interface {
	SendTest(*domain.SoftwarePkgBasicInfo) (int, error)
	ClosePR(int) error
}
