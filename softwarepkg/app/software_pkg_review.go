package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
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

func (s *softwarePkgService) Review(pid string, user *domain.Reviewer, reviews []domain.CheckItemReviewInfo) (err error) {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	approved, err := pkg.AddReview(&domain.UserReview{
		Reviewer: *user,
		Reviews:  reviews,
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

func (s *softwarePkgService) Reject(pid string, user *domain.Reviewer) error {
	pkg, version, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err = pkg.RejectBy(user); err != nil {
		return err
	}

	return s.repo.SaveSoftwarePkg(&pkg, version)
}
