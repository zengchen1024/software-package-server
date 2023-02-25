package repository

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type OptToFindSoftwarePkgs struct {
	Importer dp.Account
	Phase    dp.PackagePhase

	PageNum      int
	CountPerPage int
}

type SoftwarePkg interface {
	// AddSoftwarePkg adds a new pkg
	AddSoftwarePkg(*domain.SoftwarePkgBasicInfo) error

	SaveSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo, version int) error

	FindSoftwarePkgBasicInfo(pid string) (domain.SoftwarePkgBasicInfo, error)

	FindSoftwarePkg(pid string) (domain.SoftwarePkg, error)

	FindSoftwarePkgs(OptToFindSoftwarePkgs) (r []domain.SoftwarePkgBasicInfo, total int, err error)

	AddIssueComment(pid string, comment *domain.SoftwarePkgReviewComment) error
}
