package app

import (
	"errors"

	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/sensitivewords"
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
	if err = s.sensitive.CheckSensitiveWords(cmd.Content.ReviewComment()); err != nil {
		if sensitivewords.IsErrorSensitiveInfo(err) {
			code = errorSoftwarePkgCommentIllegal
		}

		return
	}

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
	cmd *CmdToTranslateReviewComment,
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
	_ = s.repo.AddTranslatedReviewComment(cmd.PkgId, &translated)

	return
}

func (s *softwarePkgService) Approve(pid string, user *domain.User) (code string, err error) {
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

func (s *softwarePkgService) Reject(pid string, user *domain.User) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if code, err = s.checkPermission(&pkg, user); err != nil {
		return
	}

	if _, err = pkg.RejectBy(user); err == nil {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) Abandon(pid string, user *domain.User) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if err = pkg.Abandon(user); err == nil {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) RerunCI(pid string, user *domain.User) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if err = pkg.RerunCI(user); err != nil {
		return
	}

	if err = s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		return
	}

	e := domain.NewSoftwarePkgAppUpdatedEvent(&pkg)
	if err = s.message.NotifyPkgToRerunCI(&e); err != nil {
		logrus.Errorf(
			"failed to notify re-running ci for pkg:%s, err:%s",
			pkg.Id, err.Error(),
		)
	} else {
		logrus.Debugf(
			"successfully to notify re-running ci for pkg:%s", pkg.Id,
		)
	}

	return
}

func (s *softwarePkgService) checkPermission(pkg *domain.SoftwarePkgBasicInfo, user *domain.User) (
	string, error,
) {
	if s.maintainer.HasPermission(pkg, user) {
		return "", nil
	}

	return errorSoftwarePkgNoPermission, errors.New("no permission")
}
