package app

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

func (s *softwarePkgService) GetPkgReviewDetail(pid string) (SoftwarePkgReviewDTO, string, error) {
	v, _, err := s.repo.Find(pid)
	if err != nil {
		return SoftwarePkgReviewDTO{}, errorCodeForFindingPkg(err), err
	}

	comments, err := s.commentRepo.FindReviewComments(pid)
	if err != nil {
		return SoftwarePkgReviewDTO{}, "", err
	}

	return toSoftwarePkgReviewDTO(&v, comments), "", nil
}

func (s *softwarePkgService) Review(pid string, user *domain.Reviewer, reviews []domain.CheckItemReviewInfo) error {
	pkg, version, err := s.repo.Find(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	err = pkg.AddReview(&domain.UserReview{
		Reviewer: *user,
		Reviews:  reviews,
	})
	if err != nil {
		return err
	}

	return s.repo.Save(&pkg, version)
}

func (s *softwarePkgService) Reject(pid string, user *domain.Reviewer) error {
	pkg, version, err := s.repo.FindAndIgnoreReview(pid)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err = pkg.RejectBy(user); err != nil {
		return err
	}

	return s.repo.SaveAndIgnoreReview(&pkg, version)
}
