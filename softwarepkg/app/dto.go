package app

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/utils"
)

type CmdToApplyNewSoftwarePkg struct {
	Sig      dp.ImportingPkgSig
	Spec     dp.URL
	SRPM     dp.URL
	Repo     domain.SoftwarePkgRepo
	Basic    domain.SoftwarePkgBasicInfo
	Importer domain.User
}

type CmdToUpdateSoftwarePkgApplication struct {
	PkgId    string
	Importer domain.User

	domain.SoftwarePkgUpdateInfo
}

type CmdToListPkgs = repository.OptToFindSoftwarePkgs

type CmdToWriteSoftwarePkgReviewComment struct {
	PkgId   string
	Author  dp.Account
	Content dp.ReviewComment
}

type CmdToAbandonPkg struct {
	PkgId    string
	Comment  dp.ReviewComment
	Importer dp.Account
}

type NewSoftwarePkgDTO struct {
	Id string `json:"id"`
}

// SoftwarePkgDTO
type SoftwarePkgDTO struct {
	Id        string `json:"id"`
	Sig       string `json:"sig"`
	Phase     string `json:"phase"`
	PkgName   string `json:"pkg_name"`
	PkgDesc   string `json:"desc"`
	Platform  string `json:"platform"`
	RepoLink  string `json:"repo_link"`
	CIStatus  string `json:"ci_status"`
	Importer  string `json:"importer"`
	AppliedAt string `json:"applied_at"`
}

func toSoftwarePkgDTO(v *repository.SoftwarePkgInfo) SoftwarePkgDTO {
	return SoftwarePkgDTO{
		Id:        v.Id,
		Sig:       v.Sig.ImportingPkgSig(),
		Phase:     v.Phase.PackagePhase(),
		PkgName:   v.PkgName.PackageName(),
		PkgDesc:   v.PkgDesc.PackageDesc(),
		Platform:  v.Platform.PackagePlatform(),
		RepoLink:  v.Platform.RepoLink(v.PkgName),
		CIStatus:  v.CIStatus.PackageCIStatus(),
		Importer:  v.Importer.Account(),
		AppliedAt: utils.ToDate(v.AppliedAt),
	}
}

func toSoftwarePkgDTOs(v []repository.SoftwarePkgInfo) (r []SoftwarePkgDTO) {
	if n := len(v); n > 0 {
		r = make([]SoftwarePkgDTO, n)
		for i := range v {
			r[i] = toSoftwarePkgDTO(&v[i])
		}
	}

	return
}

// SoftwarePkgApplicationDTO
type SoftwarePkgApplicationDTO struct {
	Id string `json:"id"`

	Phase     string `json:"phase"`
	CIStatus  string `json:"ci_status"`
	Importer  string `json:"importer"`
	AppliedAt string `json:"applied_at"`

	Purpose  string `json:"purpose"`
	PkgName  string `json:"name"`
	PkgDesc  string `json:"desc"`
	Upstream string `json:"upstream"`

	Spec string `json:"spec"`
	SRPM string `json:"srpm"`
	// local path , failed reason

	Sig        string   `json:"sig"`
	RepoLink   string   `json:"repo_link"`
	Committers []string `json:"committers"`
}

func toSoftwarePkgApplicationDTO(v *domain.SoftwarePkg) SoftwarePkgApplicationDTO {
	return SoftwarePkgApplicationDTO{
		Spec:     v.Code.Spec.Src.URL(),
		Upstream: v.Basic.Upstream.URL(),
		SRPM:     v.Code.SRPM.Src.URL(),
		PkgDesc:  v.Basic.Desc.PackageDesc(),
		Sig:      v.Sig.ImportingPkgSig(),
		Purpose:  v.Basic.Purpose.PurposeToImportPkg(),
	}
}

// SoftwarePkgReviewCommentDTO
type SoftwarePkgReviewCommentDTO struct {
	Id            string `json:"id"`
	Author        string `json:"author"`
	Content       string `json:"content"`
	CreatedAt     string `json:"created_at"`
	SinceCreation int64  `json:"since_creation"`
}

func toSoftwarePkgReviewCommentDTO(v *domain.SoftwarePkgReviewComment) SoftwarePkgReviewCommentDTO {
	return SoftwarePkgReviewCommentDTO{
		Id:            v.Id,
		Author:        v.Author.Account(),
		Content:       v.Content.ReviewComment(),
		CreatedAt:     utils.ToDateTime(v.CreatedAt),
		SinceCreation: utils.Now() - v.CreatedAt,
	}
}

func toSoftwarePkgReviewCommentDTOs(v []domain.SoftwarePkgReviewComment) (r []SoftwarePkgReviewCommentDTO) {
	if n := len(v); n > 0 {
		r = make([]SoftwarePkgReviewCommentDTO, n)
		for i := range v {
			r[i] = toSoftwarePkgReviewCommentDTO(&v[i])
		}
	}

	return
}

// SoftwarePkgOperationLogDTO
type SoftwarePkgOperationLogDTO struct {
	User   string `json:"user"`
	Time   string `json:"time"`
	Action string `json:"action"`
}

func toSoftwarePkgOperationLogDTO(v *domain.SoftwarePkgOperationLog) SoftwarePkgOperationLogDTO {
	return SoftwarePkgOperationLogDTO{
		User:   v.User.Account(),
		Time:   utils.ToDateTime(v.Time),
		Action: v.Action.PackageOperationLogAction(),
	}
}

func toSoftwarePkgOperationLogDTOs(v []domain.SoftwarePkgOperationLog) (r []SoftwarePkgOperationLogDTO) {
	if n := len(v); n > 0 {
		r = make([]SoftwarePkgOperationLogDTO, n)
		for i := range v {
			r[i] = toSoftwarePkgOperationLogDTO(&v[i])
		}
	}

	return
}

// SoftwarePkgApproverDTO
type SoftwarePkgApproverDTO struct {
	Account string `json:"account"`
	IsTC    bool   `json:"is_tc"`
}

// SoftwarePkgReviewDTO
type SoftwarePkgReviewDTO struct {
	Logs        []SoftwarePkgOperationLogDTO  `json:"logs"`
	Comments    []SoftwarePkgReviewCommentDTO `json:"comments"`
	Application SoftwarePkgApplicationDTO     `json:"application"`
}

func toSoftwarePkgReviewDTO(v *domain.SoftwarePkg, comments []domain.SoftwarePkgReviewComment) SoftwarePkgReviewDTO {
	return SoftwarePkgReviewDTO{
		Logs:        toSoftwarePkgOperationLogDTOs(v.Logs),
		Comments:    toSoftwarePkgReviewCommentDTOs(comments),
		Application: toSoftwarePkgApplicationDTO(v),
	}
}

// SoftwarePkgsDTO
type SoftwarePkgsDTO struct {
	Pkgs  []SoftwarePkgDTO `json:"pkgs"`
	Total int              `json:"total"`
}

func toSoftwarePkgsDTO(v []repository.SoftwarePkgInfo, total int) SoftwarePkgsDTO {
	return SoftwarePkgsDTO{
		Pkgs:  toSoftwarePkgDTOs(v),
		Total: total,
	}
}

// CmdToTranslateReviewComment
type CmdToTranslateReviewComment = repository.TranslatedReviewCommentIndex

type TranslatedReveiwCommentDTO struct {
	Content string `json:"content"`
}
