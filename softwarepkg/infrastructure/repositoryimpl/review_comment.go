package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

type reviewCommentTable struct {
	cli dbClient
}

func (s reviewCommentTable) findSoftwarePkgReviews(pid string) (
	[]domain.SoftwarePkgReviewComment, error,
) {
	var dos []SoftwarePkgReviewCommentDO

	err := s.cli.GetRecords(
		&SoftwarePkgReviewCommentDO{PkgId: pid},
		&dos,
		postgresql.Pagination{},
		[]postgresql.SortByColumn{
			{Column: createdAt, Ascend: true},
		},
	)
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	v := make([]domain.SoftwarePkgReviewComment, len(dos))
	for i, do := range dos {
		if v[i], err = do.toSoftwarePkgReviewComment(); err != nil {
			return nil, err
		}
	}

	return v, nil
}
