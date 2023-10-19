package domain

import (
	"sort"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// SoftwarePkgReview
type SoftwarePkgReview struct {
	Items   []CheckItem
	Reviews []UserReview
}

func (r *SoftwarePkgReview) add(ur *UserReview) {
	r.Reviews = append(r.Reviews, *ur)
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

	if len(rs) > 0 {
		sort.Slice(rs, func(i, j int) bool {
			oi := item.isOwner(rs[i].Role)
			oj := item.isOwner(rs[j].Role)

			return oi == oj || oi
		})
	}

	rf.Reviews = rs

	return
}

// Reviewer
type Reviewer struct {
	User dp.Account
	Role []string
}

func (r *Reviewer) isTC() bool {
	return false // TODO
}

// UserReview
type UserReview struct {
	Reviewer

	Reviews []CheckItemReviewInfo
}

func (r *UserReview) CheckItemReview(item *CheckItem) (info UserCheckItemReview, exist bool) {
	for i := range r.Reviews {
		if v := &r.Reviews[i]; v.Index == item.Index {
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
		if v := &r.Reviews[i]; r.Item.isOwner(v.Role) {
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
	Index int
	Pass  bool
	Desc  string
}

// CheckItem
type CheckItem struct {
	Index  int
	Item   string
	Desc   string
	Owners []string // TC, Sig maintainer, Committer
}

func (item *CheckItem) isOwner(Role []string) bool {
	for _, role := range Role {
		for _, owner := range item.Owners {
			if role == owner {
				return true
			}
		}
	}

	return false
}
