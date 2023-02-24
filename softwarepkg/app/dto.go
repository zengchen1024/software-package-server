package app

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type CmdToApplyNewSoftwarePkg = domain.Application
type CmdToListPkgs = repository.OptToFindSoftwarePkgs

type SoftwarePkgBasicInfoDTO struct {
	Id        string `json:"id"`
	Importer  string `json:"importer"`
	PkgName   string `json:"pkg_name"`
	Status    string `json:"status"`
	AppliedAt string `json:"applied_at"`
}

type SoftwarePkgApplicationDTO struct {
}

type SoftwarePkgIssueCommentDTO struct {
}

type SoftwarePkgIssueInfoDTO struct {
	Application SoftwarePkgApplicationDTO    `json:"application"`
	Comments    []SoftwarePkgIssueCommentDTO `json:"comments"`
	ApprovedBy  []string                     `json:"approved_by"`
	RejectedBy  []string                     `json:"rejected_by"`
}

type SoftwarePkgIssueDTO struct {
	SoftwarePkgBasicInfoDTO
	SoftwarePkgIssueInfoDTO
}

func toSoftwarePkgIssueDTO(v *domain.SoftwarePkgIssue) SoftwarePkgIssueDTO {
	return SoftwarePkgIssueDTO{}
}

type SoftwarePkgsDTO struct {
	Pkgs  []SoftwarePkgBasicInfoDTO `json:"pkgs"`
	Total int                       `json:"total"`
}
