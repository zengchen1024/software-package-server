package controller

import (
	"errors"

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
	Name     string `json:"name"        binding:"required"`
	Desc     string `json:"desc"        binding:"required"`
	Purpose  string `json:"purpose"     binding:"required"`
	Upstream string `json:"upstream"    binding:"required"`

	Spec string `json:"spec"  binding:"required"`
	SRPM string `json:"srpm"  binding:"required"`

	Sig        string   `json:"sig"          binding:"required"`
	RepoLink   string   `json:"repo_link"    binding:"required"`
	Committers []string `json:"committers"   binding:"required"`
}

func (req *softwarePkgRequest) toBasic() (basic domain.SoftwarePkgBasicInfo, err error) {
	basic.Name, err = dp.NewPackageName(req.Name)
	if err != nil {
		return
	}

	basic.Desc, err = dp.NewPackageDesc(req.Desc)
	if err != nil {
		return
	}

	basic.Purpose, err = dp.NewPurposeToImportPkg(req.Purpose)
	if err != nil {
		return
	}

	basic.Upstream, err = dp.NewURL(req.Upstream)

	return
}

func (req *softwarePkgRequest) toRepo(importer *domain.User, ua useradapter.UserAdapter) (
	repo domain.SoftwarePkgRepo, err error,
) {
	repo.Platform, err = dp.NewPackagePlatformByRepoLink(req.RepoLink)
	if err != nil {
		return
	}

	repo.Committers, err = toCommitters(importer, ua, req.Committers, repo.Platform)

	return
}

func (req *softwarePkgRequest) toCmd(importer *domain.User, ua useradapter.UserAdapter) (
	cmd app.CmdToApplyNewSoftwarePkg, err error,
) {
	cmd.Importer = *importer

	if cmd.Basic, err = req.toBasic(); err != nil {
		return
	}

	cmd.Spec, err = dp.NewURL(req.Spec)
	if err != nil {
		return
	}

	cmd.SRPM, err = dp.NewURL(req.SRPM)
	if err != nil {
		return
	}

	cmd.Sig, err = dp.NewImportingPkgSig(req.Sig)
	if err != nil {
		return
	}

	cmd.Repo, err = req.toRepo(importer, ua)

	return
}

type reqToUpdateSoftwarePkgApplication struct {
	Desc     string `json:"desc"`
	Purpose  string `json:"purpose"`
	Upstream string `json:"upstream"`

	Spec string `json:"spec"`
	SRPM string `json:"srpm"`

	Sig        string   `json:"sig"`
	RepoLink   string   `json:"repo_link"`
	Committers []string `json:"committers"`
}

func (req *reqToUpdateSoftwarePkgApplication) toCmd(
	pkgId string, importer *domain.User, ua useradapter.UserAdapter,
) (
	cmd app.CmdToUpdateSoftwarePkgApplication, err error,
) {
	cmd.PkgId = pkgId
	cmd.Importer = *importer

	if req.Desc != "" {
		if cmd.Desc, err = dp.NewPackageDesc(req.Desc); err != nil {
			return
		}
	}

	if req.Purpose != "" {
		cmd.Purpose, err = dp.NewPurposeToImportPkg(req.Purpose)
		if err != nil {
			return
		}
	}

	if req.Upstream != "" {
		if cmd.Upstream, err = dp.NewURL(req.Upstream); err != nil {
			return
		}
	}

	if req.Spec != "" {
		if cmd.Spec, err = dp.NewURL(req.Spec); err != nil {
			return
		}
	}

	if req.SRPM != "" {
		if cmd.SRPM, err = dp.NewURL(req.SRPM); err != nil {
			return
		}
	}

	if req.Sig != "" {
		if cmd.Sig, err = dp.NewImportingPkgSig(req.Sig); err != nil {
			return
		}
	}

	if req.RepoLink == "" && len(req.Committers) == 0 {
		return
	}

	if !(req.RepoLink != "" && len(req.Committers) != 0) {
		err = errors.New("repo_link and committers must be set at same time")

		return
	}

	if cmd.Repo.Platform, err = dp.NewPackagePlatformByRepoLink(req.RepoLink); err != nil {
		return
	}

	cmd.Repo.Committers, err = toCommitters(importer, ua, req.Committers, cmd.Repo.Platform)

	return
}

type softwarePkgListQuery struct {
	Phase        string `json:"phase"          form:"phase"`
	PkgName      string `json:"pkg_name"       form:"pkg_name"`
	Importer     string `json:"importer"       form:"importer"`
	Platform     string `json:"platform"       form:"platform"`
	LastId       string `json:"last_id"        form:"last_id"`
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

	if s.LastId != "" && s.PageNum > 0 {
		err = errors.New("it can't set last_id and page_num at same time")

		return
	}

	if s.LastId == "" {
		if s.PageNum <= 0 || s.PageNum > config.MaxPageNum {
			s.PageNum = pageNum
		}
		pkg.PageNum = s.PageNum
	}

	if s.CountPerPage <= 0 || s.CountPerPage > config.MaxCountPerPage {
		s.CountPerPage = countPerPage
	}
	pkg.CountPerPage = s.CountPerPage

	return
}

type reqToAbandonPkg struct {
	Comment string `json:"comment"`
}

func (req *reqToAbandonPkg) toCmd(pkgId string, user *domain.User) (cmd app.CmdToAbandonPkg, err error) {
	if req.Comment != "" {
		if cmd.Comment, err = dp.NewReviewComment(req.Comment); err != nil {
			return
		}
	}

	cmd.PkgId = pkgId
	cmd.Importer = user.Account

	return
}

type reviewCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

func (r reviewCommentRequest) toCmd(pkgId string, user *domain.User) (rc app.CmdToWriteSoftwarePkgReviewComment, err error) {
	if rc.Content, err = dp.NewReviewComment(r.Comment); err != nil {
		return
	}

	rc.PkgId = pkgId
	rc.Author = user.Account

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
	Id      string `json:"id"       binding:"required"`
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

func toCommitters(importer *domain.User, ua useradapter.UserAdapter, committers []string, p dp.PackagePlatform) (
	[]domain.PkgCommitter, error,
) {
	if len(committers) > 3 { // TODO config
		return nil, errors.New("too many committers")
	}

	platform := p.PackagePlatform()

	platformId := importer.Id(platform)
	if platformId == "" {
		return nil, errors.New("no platform Id")
	}

	r := []domain.PkgCommitter{{Account: importer.Account, PlatformId: platformId}}

	if len(committers) == 0 {
		return r, nil
	}

	for _, c := range committers {
		if c == platformId {
			continue
		}

		u, err := ua.Find(c, platform)
		if err != nil {
			return nil, err
		}

		r = append(r, domain.PkgCommitter{Account: u.Account, PlatformId: c})
	}

	return r, nil
}
