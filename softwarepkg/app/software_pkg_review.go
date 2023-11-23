package app

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

func (s *softwarePkgService) GetReview(pid string, user *domain.User) ([]CheckItemUserReviewDTO, error) {
	pkg, _, err := s.repo.Find(pid)
	if err != nil {
		return nil, parseErrorForFindingPkg(err)
	}

	userReview := pkg.UserReview(user)

	items := pkg.CheckItems()
	r := make([]CheckItemUserReviewDTO, len(items))

	for i := range items {
		item := &items[i]

		canReview, info := userReview.CheckItemReview(item)

		r[i] = CheckItemUserReviewDTO{
			Id:        item.Id,
			Name:      item.Name,
			Desc:      item.Desc,
			Owner:     item.OwnerDesc(&pkg),
			CanReview: canReview,
		}

		if canReview && info != nil {
			r[i].Pass = info.Pass
			r[i].Reviewed = true
			r[i].Comment = info.Comment
		}
	}

	return r, nil
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
