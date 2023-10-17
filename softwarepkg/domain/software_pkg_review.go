package domain

import (
	"sort"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// CheckItemReviewInfos
type CheckItemReviewInfos struct {
	Item  *CheckItem
	Infos []CheckItemReviewInfo
}

func (r *CheckItemReviewInfos) Result() dp.CheckItemResult {
	if len(r.Infos) == 0 {
		return dp.CheckItemNoIdea
	}

	pass := false

	for i := range r.Infos {
		if v := &r.Infos[i]; r.Item.isOwner(v.Role) {
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

func (r *SoftwarePkgReview) CheckItemReview(item *CheckItem) (rf CheckItemReviewInfos) {
	rf.Item = item

	if len(r.Reviews) == 0 {
		return
	}

	infos := make([]CheckItemReviewInfo, 0, len(r.Reviews))

	for i := range r.Reviews {
		if v, exist := r.Reviews[i].CheckItemReview(item); exist {
			infos = append(infos, v)
		}
	}

	if len(infos) > 0 {
		sort.Slice(infos, func(i, j int) bool {
			oi := item.isOwner(infos[i].Role)
			oj := item.isOwner(infos[j].Role)

			return oi == oj || oi
		})
	}

	rf.Infos = infos

	return
}

// CheckItemReviewInfo
type CheckItemReviewInfo struct {
	*Reviewer
	*CheckItemReview
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

	Items []CheckItemReview
}

func (r *UserReview) CheckItemReview(item *CheckItem) (info CheckItemReviewInfo, exist bool) {
	for i := range r.Items {
		if v := &r.Items[i]; v.Index == item.Index {
			exist = true

			info.Reviewer = &r.Reviewer
			info.CheckItemReview = v

			return
		}
	}

	return
}

// CheckItemReview
type CheckItemReview struct {
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
