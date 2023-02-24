package app

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type CmdToApplyNewSoftwarePkg = domain.SoftwarePkgApplication

type CmdToListPkgs = repository.OptToFindSoftwarePkgs

// SoftwarePkgBasicInfoDTO
type SoftwarePkgBasicInfoDTO struct {
	Id        string `json:"id"`
	Importer  string `json:"importer"`
	PkgName   string `json:"pkg_name"`
	Status    string `json:"status"`
	AppliedAt string `json:"applied_at"`
	RepoLink  string `json:"repo_link"`
}

func toSoftwarePkgBasicInfoDTO(v *domain.SoftwarePkgBasicInfo) SoftwarePkgBasicInfoDTO {
	dto := SoftwarePkgBasicInfoDTO{
		Id:        v.Id,
		Importer:  v.Importer.Account(),
		PkgName:   v.PkgName.PackageName(),
		Status:    v.Status.PackageStatus(),
		AppliedAt: "",
	}

	if v.RepoLink != nil {
		dto.RepoLink = v.RepoLink.URL()
	}

	return dto
}

func toSoftwarePkgsBasicInfoDTO(v []domain.SoftwarePkgBasicInfo) (r []SoftwarePkgBasicInfoDTO) {
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
	PackageName       string `json:"name"`
	PackageDesc       string `json:"desc"`
	SourceCodeLink    string `json:"source_code"`
	PackagePlatform   string `json:"platform"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason"`
	SourceCodeLicense string `json:"license"`
}

func toSoftwarePkgApplicationDTO(v *domain.SoftwarePkgApplication) SoftwarePkgApplicationDTO {
	return SoftwarePkgApplicationDTO{
		PackageName:       v.PackageName.PackageName(),
		PackageDesc:       v.PackageDesc.PackageDesc(),
		SourceCodeLink:    v.SourceCode.Address.URL(),
		PackagePlatform:   v.PackagePlatform.PackagePlatform(),
		ImportingPkgSig:   v.ImportingPkgSig.ImportingPkgSig(),
		ReasonToImportPkg: v.ReasonToImportPkg.ReasonToImportPkg(),
		SourceCodeLicense: v.SourceCode.License.License(),
	}
}

// SoftwarePkgIssueCommentDTO
type SoftwarePkgIssueCommentDTO struct {
}

func toSoftwarePkgIssueCommentDTO(v *domain.SoftwarePkgIssueComment) SoftwarePkgIssueCommentDTO {
	return SoftwarePkgIssueCommentDTO{}
}

func toSoftwarePkgIssueCommentsDTO(v []domain.SoftwarePkgIssueComment) (r []SoftwarePkgIssueCommentDTO) {
	if n := len(v); n > 0 {
		r = make([]SoftwarePkgIssueCommentDTO, n)
		for i := range v {
			r[i] = toSoftwarePkgIssueCommentDTO(&v[i])
		}
	}

	return
}

// SoftwarePkgIssueInfoDTO
type SoftwarePkgIssueInfoDTO struct {
	Application SoftwarePkgApplicationDTO    `json:"application"`
	Comments    []SoftwarePkgIssueCommentDTO `json:"comments"`
	ApprovedBy  []string                     `json:"approved_by"`
	RejectedBy  []string                     `json:"rejected_by"`
}

func toSoftwarePkgIssueInfoDTO(v *domain.SoftwarePkgIssueInfo) SoftwarePkgIssueInfoDTO {
	return SoftwarePkgIssueInfoDTO{
		Application: toSoftwarePkgApplicationDTO(&v.Application),
		Comments:    toSoftwarePkgIssueCommentsDTO(v.Comments),
		ApprovedBy:  toAccounts(v.ApprovedBy),
		RejectedBy:  toAccounts(v.RejectedBy),
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

// SoftwarePkgIssueDTO
type SoftwarePkgIssueDTO struct {
	SoftwarePkgBasicInfoDTO
	SoftwarePkgIssueInfoDTO
}

func toSoftwarePkgIssueDTO(v *repository.SoftwarePkgIssue) SoftwarePkgIssueDTO {
	return SoftwarePkgIssueDTO{
		SoftwarePkgBasicInfoDTO: toSoftwarePkgBasicInfoDTO(&v.SoftwarePkgBasicInfo),
		SoftwarePkgIssueInfoDTO: toSoftwarePkgIssueInfoDTO(&v.SoftwarePkgIssueInfo),
	}
}

// SoftwarePkgsDTO
type SoftwarePkgsDTO struct {
	Pkgs  []SoftwarePkgBasicInfoDTO `json:"pkgs"`
	Total int                       `json:"total"`
}

func toSoftwarePkgsDTO(v []domain.SoftwarePkgBasicInfo, total int) SoftwarePkgsDTO {
	return SoftwarePkgsDTO{
		Pkgs:  toSoftwarePkgsBasicInfoDTO(v),
		Total: total,
	}
}
