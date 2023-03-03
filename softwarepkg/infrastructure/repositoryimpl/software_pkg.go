package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	softwarePkgBasic

	reviewComment
}

func NewSoftwarePkg(cfg *Config) repository.SoftwarePkg {
	return softwarePkgImpl{
		softwarePkgBasic: softwarePkgBasic{
			postgresql.NewDBTable(cfg.Table.SoftwarePkg),
		},
		reviewComment: reviewComment{
			postgresql.NewDBTable(cfg.Table.ReviewComment),
		},
	}
}

func (impl softwarePkgImpl) FindSoftwarePkg(pid string) (
	pkg domain.SoftwarePkg, version int, err error,
) {
	pkg.SoftwarePkgBasicInfo, version, err = impl.softwarePkgBasic.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	pkg.Comments, err = impl.findSoftwarePkgReviews(pid)

	return
}
