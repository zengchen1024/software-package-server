package app

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/utils"
)

type CmdToApplyNewSoftwarePkg struct {
	PkgName dp.PackageName

	domain.SoftwarePkgApplication
}

type CmdToListPkgs = repository.OptToFindSoftwarePkgs

type CmdToWriteSoftwarePkgReviewComment struct {
	Author  dp.Account
	Content dp.ReviewComment
}

// SoftwarePkgBasicInfoDTO
type SoftwarePkgBasicInfoDTO struct {
	Id        string `json:"id"`
	Importer  string `json:"importer"`
	PkgName   string `json:"pkg_name"`
	Phase     string `json:"phase"`
	AppliedAt string `json:"applied_at"`
	RepoLink  string `json:"repo_link"`
}

func toSoftwarePkgBasicInfoDTO(v *domain.SoftwarePkgBasicInfo) SoftwarePkgBasicInfoDTO {
	dto := SoftwarePkgBasicInfoDTO{
		Id:        v.Id,
		Importer:  v.Importer.Account(),
		PkgName:   v.PkgName.PackageName(),
		Phase:     v.Phase.PackagePhase(),
		AppliedAt: utils.ToDate(v.AppliedAt),
	}

	if v.RepoLink != nil {
		dto.RepoLink = v.RepoLink.URL()
	}

	return dto
}

func toSoftwarePkgBasicInfoDTOs(v []domain.SoftwarePkgBasicInfo) (r []SoftwarePkgBasicInfoDTO) {
	if n := len(v); n > 0 {
		r = make([]SoftwarePkgBasicInfoDTO, n)
		for i := range v {
			r[i] = toSoftwarePkgBasicInfoDTO(&v[i])
		}
	}

	return
}

// SoftwarePkgApplicationDTO
type SoftwarePkgApplicationDTO struct {
	PackageDesc       string `json:"desc"`
	SourceCodeLink    string `json:"source_code"`
	PackagePlatform   string `json:"platform"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason"`
	SourceCodeLicense string `json:"license"`
}

func toSoftwarePkgApplicationDTO(v *domain.SoftwarePkgApplication) SoftwarePkgApplicationDTO {
	return SoftwarePkgApplicationDTO{
		PackageDesc:       v.PackageDesc.PackageDesc(),
		SourceCodeLink:    v.SourceCode.Address.URL(),
		PackagePlatform:   v.PackagePlatform.PackagePlatform(),
		ImportingPkgSig:   v.ImportingPkgSig.ImportingPkgSig(),
		ReasonToImportPkg: v.ReasonToImportPkg.ReasonToImportPkg(),
		SourceCodeLicense: v.SourceCode.License.License(),
	}
}

// SoftwarePkgReviewCommentDTO
type SoftwarePkgReviewCommentDTO struct {
	Id        string `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func toSoftwarePkgReviewCommentDTO(v *domain.SoftwarePkgReviewComment) SoftwarePkgReviewCommentDTO {
	return SoftwarePkgReviewCommentDTO{
		Id:        v.Id,
		Author:    v.Author.Account(),
		Content:   v.Content.ReviewComment(),
		CreatedAt: utils.ToDateTime(v.CreatedAt),
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

// SoftwarePkgReviewDTO
type SoftwarePkgReviewDTO struct {
	SoftwarePkgBasicInfoDTO

	ApprovedBy  []string                      `json:"approved_by"`
	RejectedBy  []string                      `json:"rejected_by"`
	Comments    []SoftwarePkgReviewCommentDTO `json:"comments"`
	Application SoftwarePkgApplicationDTO     `json:"application"`
}

func toSoftwarePkgReviewDTO(v *domain.SoftwarePkg) SoftwarePkgReviewDTO {
	return SoftwarePkgReviewDTO{
		SoftwarePkgBasicInfoDTO: toSoftwarePkgBasicInfoDTO(&v.SoftwarePkgBasicInfo),
		ApprovedBy:              toAccounts(v.ApprovedBy),
		RejectedBy:              toAccounts(v.RejectedBy),
		Comments:                toSoftwarePkgReviewCommentDTOs(v.Comments),
		Application:             toSoftwarePkgApplicationDTO(&v.Application),
	}
}

func toAccounts(v []dp.Account) (r []string) {
	if n := len(v); n > 0 {
		r = make([]string, n)
		for i := range v {
			r[i] = v[i].Account()
		}
	}

	return
}

// SoftwarePkgsDTO
type SoftwarePkgsDTO struct {
	Pkgs  []SoftwarePkgBasicInfoDTO `json:"pkgs"`
	Total int                       `json:"total"`
}

func toSoftwarePkgsDTO(v []domain.SoftwarePkgBasicInfo, total int) SoftwarePkgsDTO {
	return SoftwarePkgsDTO{
		Pkgs:  toSoftwarePkgBasicInfoDTOs(v),
		Total: total,
	}
}
