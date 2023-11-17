package domain

import (
	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Reviewer struct {
	Account dp.Account
	GiteeID string
}

type maintainer interface {
	Roles(*SoftwarePkg, *Reviewer) []dp.CommunityRole
}

var maintainerInstance maintainer

func (pkg *SoftwarePkg) addReview(ur *UserReview, items []CheckItem) error {
	uri := ur.internal(pkg)

	if err := uri.validate(items); err != nil {
		return err
	}

	for i := range pkg.Reviews {
		if dp.IsSameAccount(pkg.Reviews[i].Reviewer.Account, ur.Reviewer.Account) {
			pkg.Reviews[i] = *ur

			return nil
		}
	}

	pkg.Reviews = append(pkg.Reviews, *ur)

	return nil
}

func (pkg *SoftwarePkg) doesPassReview(items []CheckItem) bool {
	reviews := make([]userReview, len(pkg.Reviews))

	for i := range pkg.Reviews {
		reviews[i] = pkg.Reviews[i].internal(pkg)
	}

	for i := range items {
		rf := checkItemReview(&items[i], reviews)

		if !dp.IsCheckItemPass(rf.Result()) {
			return false
		}
	}

	return true
}

func (pkg *SoftwarePkg) clearReview(categories []dp.PkgModificationCategory, items []CheckItem) {
	reviews := make([]userReview, len(pkg.Reviews))

	for i := range pkg.Reviews {
		reviews[i] = pkg.Reviews[i].internal(pkg)
	}

	m := map[string]bool{}
	for i := range categories {
		m[categories[i].PkgModificationCategory()] = true
	}

	for i := range items {
		if v := &items[i]; v.isCategory(m) {
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
	Reviewer

	Reviews []CheckItemReviewInfo
}

func (r *UserReview) internal(pkg *SoftwarePkg) userReview {
	return userReview{
		UserReview: r,
		roles:      maintainerInstance.Roles(pkg, &r.Reviewer),
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
		Reviewer:            &r.Reviewer,
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
	*Reviewer
	*CheckItemReviewInfo

	Roles []dp.CommunityRole
}

// CheckItemReviewInfo
type CheckItemReviewInfo struct {
	Id      string
	Pass    bool
	Comment string
}

// CheckItem
type CheckItem struct {
	Id    string
	Name  string
	Desc  string
	Owner dp.CommunityRole

	// This check item should be checked again when the relevant modifications happened.
	Categories []dp.PkgModificationCategory

	// If true, keep the review record of reviewer who is still the owner of this item
	// else, clear all the records about this item.
	// For example, the review about the item that the user aggreed to
	// to be committer of the pkg should be kept when the committers was changed.
	KeepOwnerReview bool

	// If true, only the owner can review this item else anyone can review.
	// For example, onlye sig maintainer can determine whether the sig of pkg is correct.
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

func (item *CheckItem) isCategory(categories map[string]bool) bool {
	for i := range item.Categories {
		if categories[item.Categories[i].PkgModificationCategory()] {
			return true
		}
	}

	return false
}
