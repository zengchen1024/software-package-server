package domain

import (
	"encoding/json"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var (
	NewSoftwarePkgRetestedEvent     = NewSoftwarePkgAppliedEvent
	NewSoftwarePkgInitializedEvent  = NewSoftwarePkgApprovedEvent
	NewSoftwarePkgCodeUpdatedEvent  = NewSoftwarePkgAppliedEvent
	NewSoftwarePkgCodeChangeedEvent = NewSoftwarePkgAppliedEvent
)

// softwarePkgAppliedEvent
type softwarePkgAppliedEvent struct {
	PkgId string `json:"pkg_id"`
}

func (e *softwarePkgAppliedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgAppliedEvent(pkg *SoftwarePkg) softwarePkgAppliedEvent {
	return softwarePkgAppliedEvent{
		PkgId: pkg.Id,
	}
}

func UnmarshalToSoftwarePkgAppliedEvent(data []byte) (e softwarePkgAppliedEvent, err error) {
	err = json.Unmarshal(data, &e)

	return
}

// softwarePkgApprovedEvent
type softwarePkgApprovedEvent struct {
	Importer          string `json:"importer"`
	PkgId             string `json:"pkg_id"`
	PkgName           string `json:"pkg_name"`
	PkgDesc           string `json:"pkg_desc"`
	SpecURL           string `json:"spec_url"`
	Upstream          string `json:"upstream"`
	SrcRPMURL         string `json:"src_rpm_url"`
	Platform          string `json:"platform"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason_to_import"`
	CIPRNum           int    `json:"ci_pr_num"`
}

func (e *softwarePkgApprovedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgApprovedEvent(pkg *SoftwarePkg) softwarePkgApprovedEvent {
	basic := &pkg.Basic
	code := &pkg.Code

	return softwarePkgApprovedEvent{
		Importer:          pkg.Importer.Account(),
		PkgId:             pkg.Id,
		PkgName:           basic.Name.PackageName(),
		PkgDesc:           basic.Desc.PackageDesc(),
		SpecURL:           code.Spec.Src.URL(),
		Upstream:          basic.Upstream.URL(),
		SrcRPMURL:         code.SRPM.Src.URL(),
		CIPRNum:           pkg.CI.Id,
		Platform:          pkg.Repo.Platform.PackagePlatform(),
		ImportingPkgSig:   pkg.Sig.ImportingPkgSig(),
		ReasonToImportPkg: basic.Reason.ReasonToImportPkg(),
	}
}

// softwarePkgAlreadyExistedEvent
type softwarePkgAlreadyExistedEvent struct {
	PkgName string `json:"pkg_name"`
}

func (e *softwarePkgAlreadyExistedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgAlreadyExistEvent(pkg dp.PackageName) softwarePkgAlreadyExistedEvent {
	return softwarePkgAlreadyExistedEvent{
		PkgName: pkg.PackageName(),
	}
}

func UnmarshalToSoftwarePkgAlreadyExistEvent(data []byte) (
	e softwarePkgAlreadyExistedEvent, err error,
) {
	err = json.Unmarshal(data, &e)

	return
}
