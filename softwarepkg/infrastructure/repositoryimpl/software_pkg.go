package repositoryimpl

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
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

func (s softwarePkgImpl) FindSoftwarePkgs(pkgs repository.OptToFindSoftwarePkgs) (
	r []domain.SoftwarePkgBasicInfo, total int, err error,
) {
	var filter SoftwarePkgDO
	if pkgs.Importer != nil {
		filter.ImportUser = pkgs.Importer.Account()
	}

	if pkgs.Phase != nil {
		filter.Phase = pkgs.Phase.PackagePhase()
	}

	if total, err = s.cli.Counts(&filter); err != nil || total == 0 {
		return
	}

	var result []SoftwarePkgDO

	err = s.cli.GetTableRecords(
		&filter, &result,
		postgresql.Pagination{
			PageNum:      pkgs.PageNum,
			CountPerPage: pkgs.CountPerPage,
		},
		[]postgresql.SortByColumn{
			{Column: applyTime},
		},
	)
	if err != nil || len(result) == 0 {
		return
	}

	r = make([]domain.SoftwarePkgBasicInfo, len(result))
	for i := range result {
		if r[i], err = result[i].toSoftwarePkgSummary(); err != nil {
			return
		}
	}

	return
}

func (s softwarePkgImpl) AddReviewComment(pid string, comment *domain.SoftwarePkgReviewComment) error {
	//TODO implement me
	return nil
}

func (s softwarePkgImpl) AddSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo) error {
	v := s.toSoftwarePkgDO(pkg)
	filter := SoftwarePkgDO{PackageName: pkg.PkgName.PackageName()}
	err := s.cli.Insert(&filter, &v)
	if err != nil && s.cli.IsRowExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return err
}
