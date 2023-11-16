package app

import (
	"errors"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/sensitivewords"
)

func (s *softwarePkgService) NewReviewComment(
	pid string, cmd *CmdToWriteSoftwarePkgReviewComment,
) (code string, err error) {
	if err = s.sensitive.CheckSensitiveWords(cmd.Content.ReviewComment()); err != nil {
		if sensitivewords.IsErrorSensitiveInfo(err) {
			code = errorSoftwarePkgCommentIllegal
		}

		return
	}

	pkg, _, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		code = errorCodeForFindingPkg(err)

		return
	}

	if !pkg.CanAddReviewComment() {
		code = errorSoftwarePkgCannotComment
		err = errors.New("can't comment now")

		return
	}

	// TODO: there is a critical case that the comment can't be added now
	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	err = s.commentRepo.AddReviewComment(pid, &comment)

	return
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
