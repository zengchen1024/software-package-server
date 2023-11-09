package app

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/sensitivewords"
)

func (s *softwarePkgService) GetPkgReviewDetail(pid string) (SoftwarePkgReviewDTO, string, error) {
	v, _, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return SoftwarePkgReviewDTO{}, errorCodeForFindingPkg(err), err
	}

	comments, err := s.commentRepo.FindReviewComments(pid)
	if err != nil {
		return SoftwarePkgReviewDTO{}, "", err
	}

	return toSoftwarePkgReviewDTO(&v, comments), "", nil
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

func (s *softwarePkgService) Review(pid string, user *domain.User, reviews []domain.CheckItemReviewInfo) (err error) {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	approved, err := pkg.AddReview(&domain.UserReview{
		User:    *user,
		Reviews: reviews,
	})
	if err != nil {
		return
	}

	if err = s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		return
	}

	if approved {
		s.notifyPkgApproved(&pkg)
	}

	return
}

func (s *softwarePkgService) notifyPkgApproved(pkg *domain.SoftwarePkg) {
	e := domain.NewSoftwarePkgApprovedEvent(pkg)
	err := s.message.NotifyPkgApproved(&e)

	msg := fmt.Sprintf(
		"notify that a pkg:%s/%s was approved", pkg.Id, pkg.Basic.Name.PackageName(),
	)
	if err == nil {
		logrus.Debugf("successfully to %s", msg)
	} else {
		logrus.Errorf("failed to %s, err:%s", msg, err.Error())
	}
}

func (s *softwarePkgService) Reject(pid string, user *domain.User) error {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err = pkg.RejectBy(user); err != nil {
		return err
	}

	return s.repo.SaveSoftwarePkg(&pkg, version)
}

func (s *softwarePkgService) Abandon(pid string, user *domain.User) error {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err = pkg.Abandon(user); err != nil {
		return err
	}

	return s.repo.SaveSoftwarePkg(&pkg, version)
}

func (s *softwarePkgService) RerunCI(pid string, user *domain.User) error {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err = pkg.RerunCI(user); err != nil {
		return err
	}

	if err = s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		return err
	}

	e := domain.NewSoftwarePkgAppUpdatedEvent(&pkg)
	if err = s.message.NotifyPkgToRerunCI(&e); err != nil {
		return err
	}

	s.addCommentToRerunCI(pid)

	return nil
}

func (s *softwarePkgService) addCommentToRerunCI(pkgId string) {
	content, _ := dp.NewReviewComment("The CI will rerun now.")
	comment := domain.NewSoftwarePkgReviewComment(s.robot, content)

	if err := s.commentRepo.AddReviewComment(pkgId, &comment); err != nil {
		logrus.Errorf(
			"failed to add a comment when reruns ci for pkg:%s, err:%s",
			pkgId, err.Error(),
		)
	}
}
