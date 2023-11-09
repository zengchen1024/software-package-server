package app

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/utils"
)

type CmdToApplyNewSoftwarePkg struct {
	Sig      dp.ImportingPkgSig
	Repo     domain.SoftwarePkgRepo
	Code     domain.SoftwarePkgCode
	Basic    domain.SoftwarePkgBasicInfo
	Importer domain.User
}

type CmdToUpdateSoftwarePkgApplication struct {
	PkgId    string
	Sig      dp.ImportingPkgSig
	Repo     domain.SoftwarePkgRepo
	Basic    domain.SoftwarePkgBasicInfo
	Importer domain.User
}

type CmdToListPkgs = repository.OptToFindSoftwarePkgs

type CmdToWriteSoftwarePkgReviewComment struct {
	Author  dp.Account
	Content dp.ReviewComment
}

type NewSoftwarePkgDTO struct {
	Id string `json:"id"`
}

// SoftwarePkgDTO
type SoftwarePkgDTO struct {
	Id        string `json:"id"`
	Importer  string `json:"importer"`
	PkgName   string `json:"pkg_name"`
	Phase     string `json:"phase"`
	CIStatus  string `json:"ci_status"`
	AppliedAt string `json:"applied_at"`
	PkgDesc   string `json:"desc"`
	Sig       string `json:"sig"`
	Platform  string `json:"platform"`
}

func toSoftwarePkgDTO(v *domain.SoftwarePkg) SoftwarePkgDTO {
	dto := SoftwarePkgDTO{
		Id:        v.Id,
		Sig:       v.Sig.ImportingPkgSig(),
		Phase:     v.Phase.PackagePhase(),
		CIStatus:  v.CI.Status().PackageCIStatus(),
		PkgDesc:   v.Basic.Desc.PackageDesc(),
		PkgName:   v.Basic.Name.PackageName(),
		Platform:  v.Repo.Platform.PackagePlatform(),
		Importer:  v.Importer.Account(),
		AppliedAt: utils.ToDate(v.AppliedAt),
	}

	return dto
}

func toSoftwarePkgDTOs(v []domain.SoftwarePkg) (r []SoftwarePkgDTO) {
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
	SpecURL           string `json:"spec_url"`
	Upstream          string `json:"upstream"`
	SrcRPMURL         string `json:"src_rpm_url"`
	PackageDesc       string `json:"desc"`
	PackagePlatform   string `json:"platform"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason"`
}

func toSoftwarePkgApplicationDTO(v *domain.SoftwarePkg) SoftwarePkgApplicationDTO {
	return SoftwarePkgApplicationDTO{
		SpecURL:           v.Code.Spec.Src.URL(),
		Upstream:          v.Basic.Upstream.URL(),
		SrcRPMURL:         v.Code.SRPM.Src.URL(),
		PackageDesc:       v.Basic.Desc.PackageDesc(),
		PackagePlatform:   v.Repo.Platform.PackagePlatform(),
		ImportingPkgSig:   v.Sig.ImportingPkgSig(),
		ReasonToImportPkg: v.Basic.Reason.ReasonToImportPkg(),
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
	SoftwarePkgDTO

	Logs        []SoftwarePkgOperationLogDTO  `json:"logs"`
	Comments    []SoftwarePkgReviewCommentDTO `json:"comments"`
	Application SoftwarePkgApplicationDTO     `json:"application"`
}

func toSoftwarePkgReviewDTO(v *domain.SoftwarePkg, comments []domain.SoftwarePkgReviewComment) SoftwarePkgReviewDTO {
	return SoftwarePkgReviewDTO{
		SoftwarePkgDTO: toSoftwarePkgDTO(v),
		Logs:           toSoftwarePkgOperationLogDTOs(v.Logs),
		Comments:       toSoftwarePkgReviewCommentDTOs(comments),
		Application:    toSoftwarePkgApplicationDTO(v),
	}
}

// SoftwarePkgsDTO
type SoftwarePkgsDTO struct {
	Pkgs  []SoftwarePkgDTO `json:"pkgs"`
	Total int              `json:"total"`
}

func toSoftwarePkgsDTO(v []domain.SoftwarePkg, total int) SoftwarePkgsDTO {
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
