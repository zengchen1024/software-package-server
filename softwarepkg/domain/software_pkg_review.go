package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Reviewer struct {
	Account dp.Account
	GiteeID string
}

type maintainer interface {
	Roles(*SoftwarePkg, *Reviewer) (tc, sigMaitainer bool)
}

var maintainerInstance maintainer

func (pkg *SoftwarePkg) addReview(ur *UserReview, items []CheckItem) error {
	if ur.isEmpty() {
		b, i := pkg.oldReviewer(&ur.Reviewer)
		if !b {
			return errors.New("invalid review")
		}

		v := pkg.Reviews
		n := len(v) - 1
		if i != n {
			v[i] = v[n]
		}
		pkg.Reviews = v[:n]

		return nil
	}

	uri := ur.internal(pkg)

	if err := uri.validate(items); err != nil {
		return err
	}

	if b, i := pkg.oldReviewer(&ur.Reviewer); b {
		pkg.Reviews[i] = *ur
	} else {
		pkg.Reviews = append(pkg.Reviews, *ur)
	}

	return nil
}

func (pkg *SoftwarePkg) oldReviewer(reviewer *Reviewer) (bool, int) {
	for i := range pkg.Reviews {
		if dp.IsSameAccount(pkg.Reviews[i].Reviewer.Account, reviewer.Account) {
			return true, i
		}
	}

	return false, 0
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

func (pkg *SoftwarePkg) clearReview(pkgms []string, items []CheckItem) {
	reviews := make([]userReview, len(pkg.Reviews))

	for i := range pkg.Reviews {
		reviews[i] = pkg.Reviews[i].internal(pkg)
	}

	ms := map[string]bool{}
	for i := range pkgms {
		ms[pkgms[i]] = true
	}

	for i := range items {
		if v := &items[i]; !v.Keep && v.needRecheck(ms) {
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

func (r *UserReview) isEmpty() bool {
	return len(r.Reviews) == 0
}

func (r *UserReview) internal(pkg *SoftwarePkg) userReview {
	tc, sigMaitainer := maintainerInstance.Roles(pkg, &r.Reviewer)

	roles := map[string]bool{}

	if tc {
		roles[dp.CommunityRoleTC.CommunityRole()] = true
	}

	if sigMaitainer {
		roles[dp.CommunityRoleSigMaintainer.CommunityRole()] = true
		roles[dp.CommunityRoleRepoMember.CommunityRole()] = true
	}

	if pkg.isCommitter(r.Account) {
		roles[dp.CommunityRoleCommitter.CommunityRole()] = true

		v := r.checkItemReviewInfo(&CheckItem{Id: r.Account.Account()})

		if v != nil && v.Pass {
			roles[dp.CommunityRoleRepoMember.CommunityRole()] = true
		}
	}

	return userReview{
		UserReview: r,
		roles:      roles,
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

	roles map[string]bool
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
		isOwner:             item.isOwnerOfItem(r.roles),
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
		if v := &r.Reviews[i]; v.isOwner {
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

	isOwner bool
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

	// If true, keep the review record of reviewer, otherwise clear all the records about
	// this item when relevant modifications happened.
	// For example, the review about the item whether the user aggreed to
	// to be committer of the pkg should be kept.
	Keep bool

	// If true, only the owner can review this item else anyone can review.
	// For example, onlye sig maintainer can determine whether the sig of pkg is correct.
	OnlyOwner bool

	// This check item should be checked again when the relevant modifications happened.
	Modifications []string
}

func (item *CheckItem) isOwnerOfItem(roles map[string]bool) bool {
	return roles != nil && roles[item.Owner.CommunityRole()]
}

func (item *CheckItem) canReview(roles map[string]bool) bool {
	return !item.OnlyOwner || item.isOwnerOfItem(roles)
}

func (item *CheckItem) needRecheck(ms map[string]bool) bool {
	for _, v := range item.Modifications {
		if ms[v] {
			return true
		}
	}

	return false
}
