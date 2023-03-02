package repositoryimpl

import (
	"github.com/google/uuid"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

// softwarePkgTable
type softwarePkgTable struct {
	cli dbClient
}

func (t softwarePkgTable) SaveSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo, version int) error {
	//TODO implement me
	return nil
}

func (t softwarePkgTable) FindSoftwarePkgBasicInfo(pid string) (
	info domain.SoftwarePkgBasicInfo, version int, err error,
) {
	v, err := uuid.Parse(pid)
	if err != nil {
		return
	}

	var do SoftwarePkgDO

	if err = t.cli.GetRecord(&SoftwarePkgDO{UUID: v}, &do); err != nil {
		if t.cli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	} else {
		version = do.Version

		info, err = do.toSoftwarePkgBasicInfo()
	}

	return
}

func (t softwarePkgTable) FindSoftwarePkgs(pkgs repository.OptToFindSoftwarePkgs) (
	r []domain.SoftwarePkgBasicInfo, total int, err error,
) {
	var filter SoftwarePkgDO
	if pkgs.Importer != nil {
		filter.ImportUser = pkgs.Importer.Account()
	}

	if pkgs.Phase != nil {
		filter.Phase = pkgs.Phase.PackagePhase()
	}

	if total, err = t.cli.Count(&filter); err != nil || total == 0 {
		return
	}

	var result []SoftwarePkgDO

	err = t.cli.GetRecords(
		&filter, &result,
		postgresql.Pagination{
			PageNum:      pkgs.PageNum,
			CountPerPage: pkgs.CountPerPage,
		},
		[]postgresql.SortByColumn{
			{Column: fieldAppliedAt},
		},
	)
	if err != nil || len(result) == 0 {
		return
	}

	r = make([]domain.SoftwarePkgBasicInfo, len(result))
	for i := range result {
		if r[i], err = result[i].toSoftwarePkgBasicInfo(); err != nil {
			return
		}
	}

	return
}

func (t softwarePkgTable) AddSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo) error {
	var do SoftwarePkgDO
	t.toSoftwarePkgDO(pkg, &do)

	err := t.cli.Insert(
		&SoftwarePkgDO{PackageName: pkg.PkgName.PackageName()},
		&do,
	)
	if err != nil && t.cli.IsRowExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return err
}
