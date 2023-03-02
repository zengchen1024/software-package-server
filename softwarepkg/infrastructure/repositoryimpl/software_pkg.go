package repositoryimpl

import (
	"github.com/google/uuid"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	softwarePkgTable

	commentTable reviewCommentTable
}

func NewSoftwarePkg(cfg *Config) repository.SoftwarePkg {
	return softwarePkgImpl{
		softwarePkgTable: softwarePkgTable{
			postgresql.NewDBTable(cfg.Table.SoftwarePkg),
		},
		commentTable: reviewCommentTable{
			postgresql.NewDBTable(cfg.Table.ReviewComment),
		},
	}
}

func (s softwarePkgTable) AddReviewComment(pid string, comment *domain.SoftwarePkgReviewComment) error {
	//TODO implement me
	return nil
}

func (s softwarePkgImpl) FindSoftwarePkg(pid string) (
	pkg domain.SoftwarePkg, version int, err error,
) {
	pkg.SoftwarePkgBasicInfo, version, err = s.softwarePkgTable.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	pkg.Comments, err = s.commentTable.findSoftwarePkgReviews(pid)

	return
}

// softwarePkgTable
type softwarePkgTable struct {
	cli dbClient
}

func (s softwarePkgTable) SaveSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo, version int) error {
	//TODO implement me
	return nil
}

func (s softwarePkgTable) FindSoftwarePkgBasicInfo(pid string) (
	info domain.SoftwarePkgBasicInfo, version int, err error,
) {
	v, err := uuid.Parse(pid)
	if err != nil {
		return
	}

	var do SoftwarePkgDO

	if err = s.cli.GetRecord(&SoftwarePkgDO{UUID: v}, &do); err != nil {
		if s.cli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	} else {
		version = do.Version

		info, err = do.toSoftwarePkgBasicInfo()
	}

	return
}

func (s softwarePkgTable) FindSoftwarePkgs(pkgs repository.OptToFindSoftwarePkgs) (
	r []domain.SoftwarePkgBasicInfo, total int, err error,
) {
	var filter SoftwarePkgDO
	if pkgs.Importer != nil {
		filter.ImportUser = pkgs.Importer.Account()
	}

	if pkgs.Phase != nil {
		filter.Phase = pkgs.Phase.PackagePhase()
	}

	if total, err = s.cli.Count(&filter); err != nil || total == 0 {
		return
	}

	var result []SoftwarePkgDO

	err = s.cli.GetRecords(
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
		if r[i], err = result[i].toSoftwarePkgBasicInfo(); err != nil {
			return
		}
	}

	return
}

func (s softwarePkgTable) AddSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo) error {
	var do SoftwarePkgDO
	s.toSoftwarePkgDO(pkg, &do)

	err := s.cli.Insert(
		&SoftwarePkgDO{PackageName: pkg.PkgName.PackageName()},
		&do,
	)
	if err != nil && s.cli.IsRowExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return err
}
