package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	softwarePkgBasic

	operationLog
}

func NewSoftwarePkg(cfg *Config) repository.SoftwarePkg {
	return softwarePkgImpl{
		softwarePkgBasic: softwarePkgBasic{
			postgresql.NewDBTable(cfg.Table.SoftwarePkgBasic),
		},
		operationLog: operationLog{
			postgresql.NewDBTable(cfg.Table.OperationLog),
		},
	}
}

func (impl softwarePkgImpl) FindSoftwarePkg(pid string) (
	pkg domain.SoftwarePkg, version int, err error,
) {
	pkg, version, err = impl.softwarePkgBasic.FindSoftwarePkg(pid)
	if err != nil {
		return
	}

	pkg.Logs, err = impl.findOperationLogs(pid)

	return
}

func (impl softwarePkgImpl) HasSoftwarePkg(pkg dp.PackageName) (bool, error) {
	filter := SoftwarePkgBasicDO{PackageName: pkg.PackageName()}

	var res SoftwarePkgBasicDO

	err := impl.softwarePkgBasic.basicDBCli.GetRecord(&filter, &res)
	if err != nil {
		if impl.softwarePkgBasic.basicDBCli.IsRowNotFound(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}
