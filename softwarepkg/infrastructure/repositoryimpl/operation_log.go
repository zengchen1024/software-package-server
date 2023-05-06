package repositoryimpl

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type operationLog struct {
	commentDBCli dbClient
}

func (t operationLog) AddOperationLog(*domain.SoftwarePkgOperationLog) error {
	return nil
}

func (t operationLog) findOperationLogs(pid string) ([]domain.SoftwarePkgOperationLog, error) {
	return nil, nil
}
