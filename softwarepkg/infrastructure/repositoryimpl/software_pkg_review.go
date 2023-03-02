package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
)

type softwarePkgReviewImpl struct {
	cli dbClient
}

func NewSoftwarePkgReview(cli dbClient) softwarePkgReviewImpl {
	return softwarePkgReviewImpl{
		cli: cli,
	}
}

func (s softwarePkgReviewImpl) FindSoftwarePkgReviews(pid string) (res []SoftwarePkgReviewDO, err error) {
	var filterPkgReview = SoftwarePkgReviewDO{SoftwarePkgUUID: pid}
	err = s.cli.GetTableRecords(
		&filterPkgReview,
		&res,
		postgresql.Pagination{},
		[]postgresql.SortByColumn{
			{Column: createAt, Ascend: true},
		},
	)

	return
}
