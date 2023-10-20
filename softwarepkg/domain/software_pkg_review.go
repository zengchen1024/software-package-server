package domain

import (
	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type maintainer interface {
	Roles(*SoftwarePkg, *User) []dp.CommunityRole
}

var maintainerInstance maintainer

// SoftwarePkgReview
type SoftwarePkgReview struct {
	Items   []CheckItem
	Reviews []UserReview
}

func (r *SoftwarePkgReview) add(pkg *SoftwarePkg, ur *UserReview) error {
	uri := ur.internal(pkg)

	if err := uri.validate(r.Items); err != nil {
		return err
	}

	for i := range r.Reviews {
		if dp.IsSameAccount(r.Reviews[i].User.Account, ur.User.Account) {
			r.Reviews[i] = *ur

			return nil
		}
	}

	r.Reviews = append(r.Reviews, *ur)

	return nil
}

func (r *SoftwarePkgReview) pass(pkg *SoftwarePkg) bool {
	reviews := make([]userReview, len(r.Reviews))

	for i := range r.Reviews {
		reviews[i] = r.Reviews[i].internal(pkg)
	}

	for i := range r.Items {
		rf := checkItemReview(&r.Items[i], reviews)

		if !dp.IsCheckItemPass(rf.Result()) {
			return false
		}
	}

	return true
}

func (r *SoftwarePkgReview) clear(pkg *SoftwarePkg, categories []dp.CheckItemCategory) {
	reviews := make([]userReview, len(r.Reviews))

	for i := range r.Reviews {
		reviews[i] = r.Reviews[i].internal(pkg)
	}

	for i := range r.Items {
		if v := &r.Items[i]; v.isCategory(categories) {
			for j := range reviews {
				reviews[j].clear(v)
			}
		}
	}
}

func checkItemReview(item *CheckItem, reviews []userReview) (rf CheckItemReview) {
	rf.Item = item

	if len(reviews) == 0 {
		return
	}

	rs := make([]UserCheckItemReview, 0, len(reviews))

	for i := range reviews {
		if v, exist := reviews[i].userCheckItemReview(item); exist {
			rs = append(rs, v)
		}
	}

	rf.Reviews = rs

	return
}

// UserReview
type UserReview struct {
	User

	Reviews []CheckItemReviewInfo
}

func (r *UserReview) internal(pkg *SoftwarePkg) userReview {
	return userReview{
		UserReview: r,
		roles:      maintainerInstance.Roles(pkg, &r.User),
	}
}

func (r *UserReview) checkItemReviewInfo(item *CheckItem) *CheckItemReviewInfo {
	for i := range r.Reviews {
		if v := &r.Reviews[i]; v.Id == item.Id {
			return v
		}
	}

	return nil
}

// userReview
type userReview struct {
	*UserReview

	roles []dp.CommunityRole
}

func (r *userReview) validate(items []CheckItem) error {
	for i := range items {
		info := r.checkItemReviewInfo(&items[i])

		if info != nil && !items[i].canReview(r.roles) {
			return allerror.NewNoPermission("not the owner")
		}
	}

	return nil
}

func (r *userReview) clear(item *CheckItem) {
	if item.KeepOwnerReview && item.isOwner(r.roles) {
		return
	}

	for i := range r.Reviews {
		if v := &r.Reviews[i]; v.Id == item.Id {
			n := len(r.Reviews) - 1
			if i != n {
				r.Reviews[i] = r.Reviews[n]
			}
			r.Reviews = r.Reviews[:n]

			return
		}
	}
}

func (r *userReview) userCheckItemReview(item *CheckItem) (UserCheckItemReview, bool) {
	info := r.checkItemReviewInfo(item)
	if info == nil {
		return UserCheckItemReview{}, false
	}

	return UserCheckItemReview{
		User:                &r.User,
		Roles:               r.roles,
		CheckItemReviewInfo: info,
	}, true
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
		if v := &r.Reviews[i]; r.Item.isOwner(v.Roles) {
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
	*User
	*CheckItemReviewInfo

	Roles []dp.CommunityRole
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

func (item *CheckItem) isOwner(roles []dp.CommunityRole) bool {
	for _, role := range roles {
		if dp.IsSameCommunityRole(role, item.Owner) {
			return true
		}
	}

	return false
}

func (item *CheckItem) canReview(roles []dp.CommunityRole) bool {
	return !item.OnlyOwnerCanReview || item.isOwner(roles)
}

func (item *CheckItem) isCategory(categories []dp.CheckItemCategory) bool {
	for i := range categories {
		if dp.IsSameCheckItemCategory(item.Category, categories[i]) {
			return true
		}
	}

	return false
}
