package pkgtester

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type SoftwarePkgInfo struct {
	PkgId      string
	SourceCode domain.SoftwarePkgSourceCode
}

type PkgTester interface {
	SendTest(*SoftwarePkgInfo) error
}
