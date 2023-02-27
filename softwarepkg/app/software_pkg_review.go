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
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(pid)
	if err != nil {
		return
	}

	if !pkg.CanAddReviewComment() {
		code = errorSoftwarePkgCannotComment
		err = errors.New("can't comment now")

		return
	}

	isCmd, isApprove := cmd.Content.ParseReviewComment()

	if isCmd {
		if code, err = s.checkPermission(&pkg, cmd.Author); err != nil {
			return
		}
	}

	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	if err = s.repo.AddReviewComment(pid, &comment); err != nil || !isCmd {
		return
	}

	var success bool
	if isApprove {
		success = s.reviewServie.ApprovePkg(&pkg, version, cmd.Author)
	} else {
		success = s.reviewServie.RejectPkg(&pkg, version, cmd.Author)
	}

	if success {
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
