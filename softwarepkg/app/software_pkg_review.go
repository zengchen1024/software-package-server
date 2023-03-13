package app

import (
	"errors"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

func (s *softwarePkgService) GetPkgReviewDetail(pid string) (SoftwarePkgReviewDTO, error) {
	v, _, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return SoftwarePkgReviewDTO{}, err
	}

	return toSoftwarePkgReviewDTO(&v), nil
}

func (s *softwarePkgService) NewReviewComment(
	pid string, cmd *CmdToWriteSoftwarePkgReviewComment,
) (code string, err error) {
	pkg, _, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if !pkg.CanAddReviewComment() {
		code = errorSoftwarePkgCannotComment
		err = errors.New("can't comment now")

		return
	}

	// TODO: there is a critical case that the comment can't be added now
	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	err = s.repo.AddReviewComment(pid, &comment)

	return
}

func (s *softwarePkgService) TranslateReviewComment(
	cmd *CmdToGetTranslatedReviewComment,
) (dto TranslatedReveiwCommentDTO, code string, err error) {
	v, err := s.repo.FindTranslatedReviewComment(cmd)
	if err == nil {
		dto.Content = v.Content

		return
	}

	if !commonrepo.IsErrorResourceNotFound(err) {
		return
	}

	// translate it
	comment, err := s.repo.FindReviewComment(cmd.PkgId, cmd.CommentId)
	if err != nil {
		if commonrepo.IsErrorResourceNotFound(err) {
			code = errorSoftwarePkgCommentNotFound
		}

		return
	}

	content, err := s.translating.Translate(
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
	_ = s.repo.AddTranslatedReviewComment(cmd.PkgId, &translated)

	return
}

func (s *softwarePkgService) Approve(pid string, user dp.Account) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if code, err = s.checkPermission(&pkg, user); err != nil {
		return
	}

	if err = s.reviewServie.ApprovePkg(&pkg, user); err == nil {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) Reject(pid string, user dp.Account) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if code, err = s.checkPermission(&pkg, user); err != nil {
		return
	}

	if err = s.reviewServie.RejectPkg(&pkg, user); err == nil {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) Abandon(pid string, user dp.Account) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if err = s.reviewServie.AbandonPkg(&pkg, user); err == nil {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) checkPermission(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) (
	code string, err error,
) {
	b, err := s.maintainer.HasPermission(pkg, user)
	if err != nil {
		return
	}

	if !b {
		code = errorSoftwarePkgNoPermission
		err = errors.New("no permission")
	}

	return
}
