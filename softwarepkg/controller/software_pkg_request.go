package controller

import (
	"errors"
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/useradapter"
)

const (
	pageNum      = 1
	countPerPage = 10
)

type softwarePkgRequest struct {
	SpecUrl         string   `json:"spec_url"        binding:"required"`
	Upstream        string   `json:"upstream"        binding:"required"`
	SrcRPMURL       string   `json:"src_rpm_url"     binding:"required"`
	PackageSig      string   `json:"sig"             binding:"required"`
	PackageName     string   `json:"pkg_name"        binding:"required"`
	PackageDesc     string   `json:"desc"            binding:"required"`
	PackageReason   string   `json:"reason"          binding:"required"`
	PackagePlatform string   `json:"platform"        binding:"required"`
	Committers      []string `json:"committers"`
}

func (s *softwarePkgRequest) toCmd(importer *domain.User, m useradapter.UserAdapter) (
	cmd app.CmdToApplyNewSoftwarePkg, err error,
) {
	cmd.Importer = *importer

	basic := &cmd.Basic

	basic.Name, err = dp.NewPackageName(s.PackageName)
	if err != nil {
		return
	}

	cmd.Code.Spec.Src, err = dp.NewURL(s.SpecUrl)
	if err != nil {
		return
	}

	basic.Upstream, err = dp.NewURL(s.Upstream)
	if err != nil {
		return
	}

	cmd.Code.SRPM.Src, err = dp.NewURL(s.SrcRPMURL)
	if err != nil {
		return
	}

	cmd.Sig, err = dp.NewImportingPkgSig(s.PackageSig)
	if err != nil {
		return
	}

	basic.Reason, err = dp.NewReasonToImportPkg(s.PackageReason)
	if err != nil {
		return
	}

	basic.Desc, err = dp.NewPackageDesc(s.PackageDesc)
	if err != nil {
		return
	}

	cmd.Repo.Platform, err = dp.NewPackagePlatform(s.PackagePlatform)
	if err != nil {
		return
	}

	users, err := toUsers(importer, m, s.Committers, cmd.Repo.Platform)
	if err != nil {
		return
	}

	cmd.Repo.Committers = toCommitters(users)

	return
}

type reqToUpdateSoftwarePkgApplication struct {
	Upstream        string `json:"upstream"        binding:"required"`
	PackageSig      string `json:"sig"             binding:"required"`
	PackageName     string `json:"pkg_name"        binding:"required"`
	PackageDesc     string `json:"desc"            binding:"required"`
	PackageReason   string `json:"reason"          binding:"required"`
	PackagePlatform string `json:"platform"        binding:"required"`
}

func (s *reqToUpdateSoftwarePkgApplication) toCmd(pkgId string, importer *domain.User) (
	cmd app.CmdToUpdateSoftwarePkgApplication, err error,
) {
	cmd.PkgId = pkgId
	cmd.Importer = *importer

	basic := &cmd.Basic

	basic.Name, err = dp.NewPackageName(s.PackageName)
	if err != nil {
		return
	}

	basic.Upstream, err = dp.NewURL(s.Upstream)
	if err != nil {
		return
	}

	cmd.Sig, err = dp.NewImportingPkgSig(s.PackageSig)
	if err != nil {
		return
	}

	basic.Reason, err = dp.NewReasonToImportPkg(s.PackageReason)
	if err != nil {
		return
	}

	basic.Desc, err = dp.NewPackageDesc(s.PackageDesc)
	if err != nil {
		return
	}

	cmd.Repo.Platform, err = dp.NewPackagePlatform(s.PackagePlatform)

	return
}

type softwarePkgListQuery struct {
	Phase        string `json:"phase"          form:"phase"`
	PkgName      string `json:"pkg_name"       form:"pkg_name"`
	Importer     string `json:"importer"       form:"importer"`
	Platform     string `json:"platform"       form:"platform"`
	PageNum      int    `json:"page_num"       form:"page_num"`
	CountPerPage int    `json:"count_per_page" form:"count_per_page"`
}

func (s softwarePkgListQuery) toCmd() (pkg app.CmdToListPkgs, err error) {
	if s.Importer != "" {
		if pkg.Importer, err = dp.NewAccount(s.Importer); err != nil {
			return
		}
	}

	if s.Phase != "" {
		if pkg.Phase, err = dp.NewPackagePhase(s.Phase); err != nil {
			return
		}
	}

	if s.Platform != "" {
		if pkg.Platform, err = dp.NewPackagePlatform(s.Platform); err != nil {
			return
		}
	}

	if s.PkgName != "" {
		if pkg.PkgName, err = dp.NewPackageName(s.PkgName); err != nil {
			return
		}
	}

	if s.PageNum > 0 {
		pkg.PageNum = s.PageNum
	} else {
		pkg.PageNum = pageNum
	}

	if s.CountPerPage > 0 {
		pkg.CountPerPage = s.CountPerPage
	} else {
		pkg.CountPerPage = countPerPage
	}

	return
}

type reviewCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

func (r reviewCommentRequest) toCmd(user *domain.User) (rc app.CmdToWriteSoftwarePkgReviewComment, err error) {
	rc.Author = user.Account

	rc.Content, err = dp.NewReviewComment(r.Comment)

	return
}

type translationCommentRequest struct {
	Language string `json:"language"`
}

func (t translationCommentRequest) toCmd(pkgId, commentId string) (cmd app.CmdToTranslateReviewComment, err error) {
	cmd.PkgId = pkgId
	cmd.CommentId = commentId
	cmd.Language, err = dp.NewLanguage(t.Language)

	return
}

type reviewRequest struct {
	Reviews []checkItemReviewInfo `json:"reviews" binding:"required"`
}

func (req *reviewRequest) toCmd() (reviews []domain.CheckItemReviewInfo, err error) {
	reviews = make([]domain.CheckItemReviewInfo, len(req.Reviews))

	for i := range req.Reviews {
		if reviews[i], err = req.Reviews[i].toInfo(); err != nil {
			return
		}
	}

	return
}

type checkItemReviewInfo struct {
	Id      int    `json:"id"       binding:"required"`
	Pass    bool   `json:"pass"`
	Comment string `json:"comment"`
}

func (req *checkItemReviewInfo) toInfo() (domain.CheckItemReviewInfo, error) {
	if !req.Pass && req.Comment == "" {
		return domain.CheckItemReviewInfo{}, errors.New(
			"lack of reasons for failure",
		)
	}

	return domain.CheckItemReviewInfo{
		Id:      req.Id,
		Pass:    req.Pass,
		Comment: req.Comment,
	}, nil
}

func toUsers(importer *domain.User, m useradapter.UserAdapter, committers []string, p dp.PackagePlatform) (
	[]domain.User, error,
) {
	if len(committers) > 3 { // TODO config
		return nil, errors.New("too many committers")
	}

	r := []domain.User{*importer}

	if len(committers) == 0 {
		return r, nil
	}

	v := map[string]bool{}
	for _, c := range committers {
		v[c] = true
	}
	if k := importer.Account.Account(); v[k] {
		delete(v, k)
	}

	if len(v) == 0 {
		return r, nil
	}

	for c := range v {
		u, err := m.Find(c)
		if err != nil {
			return nil, err
		}

		if !u.ApplyTo(p) {
			return nil, fmt.Errorf(
				"committer doesn't have account of %s",
				p.PackagePlatform(),
			)
		}

		r = append(r, u)
	}

	return r, nil
}

func toCommitters(users []domain.User) []dp.Account {
	v := make([]dp.Account, len(users))

	for i := range users {
		v[i] = users[i].Account
	}

	return v
}
