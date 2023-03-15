package repositoryimpl

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type translationComment struct {
	translationDBCli dbClient
}

func (t translationComment) FindTranslatedReviewComment(index *repository.TranslatedReviewCommentIndex) (
	r domain.SoftwarePkgTranslatedReviewComment, err error,
) {
	filter := SoftwarePkgTranslationCommentDO{
		PkgId:     index.PkgId,
		CommentId: index.CommentId,
		Language:  index.Language.Language(),
	}

	var res SoftwarePkgTranslationCommentDO
	if err = t.translationDBCli.GetRecord(&filter, &res); err != nil {
		if t.translationDBCli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotFound(err)
		}
	} else {
		r, err = res.toSoftwarePkgTranslatedReviewComment()
	}

	return
}

func (t translationComment) AddTranslatedReviewComment(
	pid string, comment *domain.SoftwarePkgTranslatedReviewComment,
) error {
	var do SoftwarePkgTranslationCommentDO
	t.toSoftwarePkgTranslationCommentDO(pid, comment, &do)

	filter := SoftwarePkgTranslationCommentDO{
		PkgId:     do.PkgId,
		CommentId: do.CommentId,
		Language:  do.Language,
	}

	return t.translationDBCli.Insert(&filter, &do)
}
