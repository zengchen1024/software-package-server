package repositoryimpl

import (
	"github.com/google/uuid"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type softwarePkgImpl struct {
	cli       dbClient
	pkgReview softwarePkgReviewImpl
}

func NewSoftwarePkg(cli dbClient, pkgReview softwarePkgReviewImpl) repository.SoftwarePkg {
	return softwarePkgImpl{cli: cli, pkgReview: pkgReview}
}

func (s softwarePkgImpl) SaveSoftwarePkg(pkg *domain.SoftwarePkgBasicInfo, version int) error {
	//TODO implement me
	return nil
}

func (s softwarePkgImpl) FindSoftwarePkgBasicInfo(pid string) (info domain.SoftwarePkgBasicInfo, version int, err error) {
	var u uuid.UUID
	if u, err = uuid.Parse(pid); err != nil {
		return
	}

	var (
		softwarePkg SoftwarePkgDO
		filterPkg   = SoftwarePkgDO{UUID: u}
	)
	if err = s.cli.GetTableRecord(&filterPkg, &softwarePkg); err != nil {
		if s.cli.IsRowNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
		return
	}

	version = softwarePkg.Version

	info, err = softwarePkg.toSoftwarePkgSummary()

	return
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

func (s softwarePkgImpl) FindSoftwarePkgReviewComments(pid string) (
	comments []domain.SoftwarePkgReviewComment, err error,
) {
	var softwarePkgReview []SoftwarePkgReviewDO
	if softwarePkgReview, err = s.pkgReview.FindSoftwarePkgReviews(pid); err != nil {
		return
	}

	comments = make([]domain.SoftwarePkgReviewComment, len(softwarePkgReview))
	for i, do := range softwarePkgReview {
		if comments[i], err = do.toSoftwarePkgReviewCommentSummary(); err != nil {
			return
		}
	}

	return
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
