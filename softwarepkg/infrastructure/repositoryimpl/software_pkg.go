package repositoryimpl

import (
	"errors"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	cli dbClient
}

func NewSoftwarePkg(cli dbClient) repository.SoftwarePkg {
	return softwarePkgImpl{cli: cli}
}

func (s softwarePkgImpl) SaveSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo, version int) error {
	//TODO implement me
	return nil
}

func (s softwarePkgImpl) FindSoftwarePkgBasicInfo(pid string) (domain.SoftwarePkgBasicInfo, int, error) {
	//TODO implement me
	return domain.SoftwarePkgBasicInfo{}, 0, nil
}

func (s softwarePkgImpl) FindSoftwarePkg(pid string) (domain.SoftwarePkg, int, error) {
	//TODO implement me
	return domain.SoftwarePkg{}, 0, nil
}

func (s softwarePkgImpl) FindSoftwarePkgs(pkgs repository.OptToFindSoftwarePkgs) (r []domain.SoftwarePkgBasicInfo, total int, err error) {
	//TODO implement me
	return nil, 0, err
}

func (s softwarePkgImpl) AddReviewComment(pid string, comment *domain.SoftwarePkgReviewComment) error {
	//TODO implement me
	return nil
}

func (s softwarePkgImpl) AddSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo) error {
	softwarePkgDO := s.toSoftwarePkgDO(pkg)

	return s.save(softwarePkgDO)
}

func (s softwarePkgImpl) save(soft *SoftwarePkgDO) error {
	filter := &SoftwarePkgDO{PackageName: soft.PackageName}
	rows, err := s.cli.Insert(filter, soft)
	if err != nil {
		return err
	}

	if rows == 0 {
		return commonrepo.NewErrorDuplicateCreating(errors.New("package exists"))
	}

	return nil
}
