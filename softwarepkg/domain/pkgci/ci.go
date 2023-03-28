package pkgci

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type SoftwarePkgInfo struct {
	PkgId      string
	SourceCode domain.SoftwarePkgSourceCode
}

type PkgCI interface {
	SendTest(*SoftwarePkgInfo) error
}
