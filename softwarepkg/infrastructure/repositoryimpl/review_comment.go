package repositoryimpl

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

type reviewComment struct {
	commentDBCli dbClient
}

func (t reviewComment) AddReviewComment(pid string, comment *domain.SoftwarePkgReviewComment) error {
	var do SoftwarePkgReviewCommentDO
	t.toSoftwarePkgReviewCommentDO(pid, comment, &do)

	filter := SoftwarePkgReviewCommentDO{Id: do.Id}

	return t.commentDBCli.Insert(&filter, &do)
}

func (t reviewComment) findSoftwarePkgReviews(pid string) (
	[]domain.SoftwarePkgReviewComment, error,
) {
	var dos []SoftwarePkgReviewCommentDO

	err := t.commentDBCli.GetRecords(
		&SoftwarePkgReviewCommentDO{PkgId: pid},
		&dos,
		postgresql.Pagination{},
		[]postgresql.SortByColumn{
			{Column: fieldCreatedAt, Ascend: true},
		},
	)
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	v := make([]domain.SoftwarePkgReviewComment, len(dos))
	for i := range dos {
		if v[i], err = dos[i].toSoftwarePkgReviewComment(); err != nil {
			return nil, err
		}
	}

	return v, nil
}
