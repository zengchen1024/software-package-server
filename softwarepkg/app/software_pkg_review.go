package app

import (
	"errors"

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

	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	err = s.repo.AddReviewComment(pid, &comment)

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

	if s.reviewServie.ApprovePkg(&pkg, version, user) {
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

	if s.reviewServie.RejectPkg(&pkg, version, user) {
		err = s.repo.SaveSoftwarePkg(&pkg, version)
	}

	return
}

func (s *softwarePkgService) Close(pid string, user dp.Account) (code string, err error) {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if !pkg.IsImporter(user) {
		if code, err = s.checkPermission(&pkg, user); err != nil {
			return
		}
	}

	if err = pkg.Close(); err == nil {
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
