package app

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

func (s *softwarePkgService) NewReviewComment(cmd *CmdToWriteSoftwarePkgReviewComment) error {
	pkg, _, err := s.repo.FindAndIgnoreReview(cmd.PkgId)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err := pkg.CanAddReviewComment(); err != nil {
		return err
	}

	// TODO: there is a critical case that the comment can't be added now
	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	return s.commentRepo.AddReviewComment(cmd.PkgId, &comment)
}

func (s *softwarePkgService) TranslateReviewComment(
	cmd *CmdToTranslateReviewComment,
) (dto TranslatedReveiwCommentDTO, code string, err error) {
	v, err := s.commentRepo.FindTranslatedReviewComment(cmd)
	if err == nil {
		dto.Content = v.Content

		return
	}

	if !commonrepo.IsErrorResourceNotFound(err) {
		return
	}

	// translate it
	comment, err := s.commentRepo.FindReviewComment(cmd.PkgId, cmd.CommentId)
	if err != nil {
		if commonrepo.IsErrorResourceNotFound(err) {
			code = errorSoftwarePkgCommentNotFound
		}

		return
	}

	content, err := s.translation.Translate(
		comment.Content.ReviewComment(), cmd.Language,
	)
	if err != nil {
		return
	}

	dto.Content = content

	// save the translated one
	translated := domain.NewSoftwarePkgTranslatedReviewComment(
		&comment, content, cmd.Language,
	)
	_ = s.commentRepo.AddTranslatedReviewComment(cmd.PkgId, &translated)

	return
}
