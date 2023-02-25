package app

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
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
		code = ""
		err = errors.New("can't comment now")

		return
	}

	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	if err = s.repo.AddReviewComment(pid, &comment); err != nil {
		return
	}

	isCmd, isApprove, isReject := cmd.Content.ParseReviewComment()
	if !isCmd {
		return
	}

	if isApprove {
		err = s.reviewServie.ApprovePkg(&pkg, version, cmd.Author)
		// TODO code no permission
	}

	if isReject {
		err = s.reviewServie.RejectPkg(&pkg, version, cmd.Author)
		// TODO code no permission
	}

	return
}
