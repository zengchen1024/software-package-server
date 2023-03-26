package domain

import "encoding/json"

// softwarePkgApprovedEvent
type softwarePkgApprovedEvent struct {
	PkgId    string `json:"pkg_id"`
	PkgName  string `json:"pkg_name"`
	PRNum    int    `json:"pr_num"`
	Platform string `json:"platform"`
}

func (e *softwarePkgApprovedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgApprovedEvent(pkg *SoftwarePkgBasicInfo) softwarePkgApprovedEvent {
	return softwarePkgApprovedEvent{
		PkgId:    pkg.Id,
		PkgName:  pkg.PkgName.PackageName(),
		Platform: pkg.Application.PackagePlatform.PackagePlatform(),
	}
}

// softwarePkgAppliedEvent
type softwarePkgAppliedEvent struct {
	Importer          string `json:"importer"`
	ImporterEmail     string `json:"importer_email"`
	PkgId             string `json:"pkg_id"`
	PkgName           string `json:"pkg_name"`
	PkgDesc           string `json:"pkg_desc"`
	SpecURL           string `json:"spec_url"`
	SrcRPMURL         string `json:"src_rpm_url"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason_to_import"`
}

func (e *softwarePkgAppliedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgAppliedEvent(
	importer *User,
	pkg *SoftwarePkgBasicInfo,
) softwarePkgAppliedEvent {
	app := &pkg.Application

	return softwarePkgAppliedEvent{
		Importer:          importer.Account.Account(),
		ImporterEmail:     importer.Email.Email(),
		PkgId:             pkg.Id,
		PkgName:           pkg.PkgName.PackageName(),
		PkgDesc:           app.PackageDesc.PackageDesc(),
		SpecURL:           app.SourceCode.SpecURL.URL(),
		SrcRPMURL:         app.SourceCode.SrcRPMURL.URL(),
		ImportingPkgSig:   app.ImportingPkgSig.ImportingPkgSig(),
		ReasonToImportPkg: app.ReasonToImportPkg.ReasonToImportPkg(),
	}
}

var NewSoftwarePkgInitializedEvent = NewSoftwarePkgApprovedEvent
