package domain

import (
	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// SoftwarePkgReview
type SoftwarePkgReview struct {
	Items   []CheckItem
	Reviews []UserReview
}

func (r *SoftwarePkgReview) add(ur *UserReview) error {
	if err := ur.validate(r.Items); err != nil {
		return err
	}

	r.Reviews = append(r.Reviews, *ur)

	return nil
}

func (r *SoftwarePkgReview) pass() bool {
	for i := range r.Items {
		if rf := r.CheckItemReview(&r.Items[i]); !dp.IsCheckItemPass(rf.Result()) {
			return false
		}
	}

	return true
}

func (r *SoftwarePkgReview) CheckItemReview(item *CheckItem) (rf CheckItemReview) {
	rf.Item = item

	if len(r.Reviews) == 0 {
		return
	}

	rs := make([]UserCheckItemReview, 0, len(r.Reviews))

	for i := range r.Reviews {
		if v, exist := r.Reviews[i].CheckItemReview(item); exist {
			rs = append(rs, v)
		}
	}

	rf.Reviews = rs

	return
}

// Reviewer
type Reviewer struct {
	User  dp.Account
	Roles []dp.CommunityRole
}

func (r *Reviewer) isTC() bool {
	for i := range r.Roles {
		if r.Roles[i].IsTC() {
			return true
		}
	}

	return false
}

// UserReview
type UserReview struct {
	Reviewer

	Reviews []CheckItemReviewInfo
}

func (r *UserReview) validate(items []CheckItem) error {
	for i := range items {
		_, exist := r.CheckItemReview(&items[i])

		if exist && !items[i].canReview(&r.Reviewer) {
			return allerror.NewNoPermission("not the owner")
		}
	}

	return nil
}

func (r *UserReview) CheckItemReview(item *CheckItem) (info UserCheckItemReview, exist bool) {
	for i := range r.Reviews {
		if v := &r.Reviews[i]; v.Id == item.Id {
			exist = true

			info.Reviewer = &r.Reviewer
			info.CheckItemReviewInfo = v

			return
		}
	}

	return
}

// CheckItemReview
type CheckItemReview struct {
	Item    *CheckItem
	Reviews []UserCheckItemReview
}

func (r *CheckItemReview) Result() dp.CheckItemReviewResult {
	if len(r.Reviews) == 0 {
		return dp.CheckItemNoIdea
	}

	pass := false

	for i := range r.Reviews {
		if v := &r.Reviews[i]; r.Item.isOwner(v.Reviewer) {
			if !v.Pass {
				return dp.CheckItemNotPass
			}

			pass = true
		}
	}

	if pass {
		return dp.CheckItemPass
	}

	return dp.CheckItemNoIdea
}

// UserCheckItemReview
type UserCheckItemReview struct {
	*Reviewer
	*CheckItemReviewInfo
}

// CheckItemReviewInfo
type CheckItemReviewInfo struct {
	Id      int
	Pass    bool
	Comment string
}

// CheckItem
type CheckItem struct {
	Id       int
	Name     string
	Desc     string
	Owner    dp.CommunityRole
	Category dp.CheckItemCategory

	// if true, keep the review record of reviewer who is still the owner of this item
	// else, clear all the records about this item
	KeepOwnerReview bool

	// if true, only the owner can review this item
	// else, anyone can review.
	OnlyOwnerCanReview bool
}

func (item *CheckItem) isOwner(reviewer *Reviewer) bool {
	for _, role := range reviewer.Roles {
		if dp.IsSameCommunityRole(role, item.Owner) {
			return true
		}
	}

	return false
}

func (item *CheckItem) canReview(reviewer *Reviewer) bool {
	return !item.OnlyOwnerCanReview || item.isOwner(reviewer)
}
