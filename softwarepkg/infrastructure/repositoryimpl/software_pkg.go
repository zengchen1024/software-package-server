package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	softwarePkgTable

	reviewCommentTable
}

func NewSoftwarePkg(cfg *Config) repository.SoftwarePkg {
	return softwarePkgImpl{
		softwarePkgTable: softwarePkgTable{
			postgresql.NewDBTable(cfg.Table.SoftwarePkg),
		},
		reviewCommentTable: reviewCommentTable{
			postgresql.NewDBTable(cfg.Table.ReviewComment),
		},
	}
}

func (impl softwarePkgImpl) FindSoftwarePkg(pid string) (
	pkg domain.SoftwarePkg, version int, err error,
) {
	pkg.SoftwarePkgBasicInfo, version, err = impl.softwarePkgTable.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	pkg.Comments, err = impl.findSoftwarePkgReviews(pid)

	return
}
