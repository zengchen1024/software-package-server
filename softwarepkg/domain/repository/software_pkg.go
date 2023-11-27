package repository

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type OptToFindSoftwarePkgs struct {
	Phase    dp.PackagePhase
	PkgName  dp.PackageName
	Platform dp.PackagePlatform
	Importer dp.Account

	Count        bool   // count the num of pkgs
	LastId       string // the id of pkg which is the last item of previous page
	PageNum      int    // it can't set both PageNum and LastId
	CountPerPage int
}

type TranslatedReviewCommentIndex struct {
	PkgId     string
	CommentId string
	Language  dp.Language
}

// SoftwarePkgInfo
type SoftwarePkgInfo struct {
	Id        string
	Sig       dp.ImportingPkgSig
	Phase     dp.PackagePhase
	PkgName   dp.PackageName
	PkgDesc   dp.PackageDesc
	Platform  dp.PackagePlatform
	CIStatus  dp.PackageCIStatus
	Importer  dp.Account
	AppliedAt int64
}

// SoftwarePkgAdapter
type SoftwarePkg interface {
	Add(*domain.SoftwarePkg) error

	Find(pid string) (domain.SoftwarePkg, int, error)
	Save(pkg *domain.SoftwarePkg, version int) error

	FindAndIgnoreReview(pid string) (domain.SoftwarePkg, int, error)
	SaveAndIgnoreReview(pkg *domain.SoftwarePkg, version int) error

	FindAll(*OptToFindSoftwarePkgs) (r []SoftwarePkgInfo, total int64, err error)
	FindAllApproved(dp.PackagePhase) ([]string, error)
}

type SoftwarePkgComment interface {
	FindReviewComments(pid string) ([]domain.SoftwarePkgReviewComment, error)

	AddReviewComment(pid string, comment *domain.SoftwarePkgReviewComment) error
	FindReviewComment(pid, commentId string) (domain.SoftwarePkgReviewComment, error)

	AddTranslatedReviewComment(pid string, comment *domain.SoftwarePkgTranslatedReviewComment) error
	FindTranslatedReviewComment(*TranslatedReviewCommentIndex) (domain.SoftwarePkgTranslatedReviewComment, error)
}
