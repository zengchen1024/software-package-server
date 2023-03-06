package domain

import (
	"encoding/json"
	"errors"
)

// softwarePkgApprovedEvent
type softwarePkgApprovedEvent struct {
	PkgId      string `json:"pkg_id"`
	PkgName    string `json:"pkg_name"`
	RelevantPR string `json:"pr"`
}

func (e *softwarePkgApprovedEvent) ToMessage() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgApprovedEvent(pkg *SoftwarePkgBasicInfo) (e softwarePkgApprovedEvent, err error) {
	if pkg.RelevantPR != nil {
		e = softwarePkgApprovedEvent{
			PkgId:      pkg.Id,
			PkgName:    pkg.PkgName.PackageName(),
			RelevantPR: pkg.RelevantPR.URL(),
		}
	} else {
		err = errors.New("missing pr")
	}

	return
}

// softwarePkgRejectedEvent
type softwarePkgRejectedEvent struct {
	RelevantPR string `json:"pr"`
}

func (e *softwarePkgRejectedEvent) ToMessage() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgRejectedEvent(pkg *SoftwarePkgBasicInfo) (e softwarePkgRejectedEvent, err error) {
	if pkg.RelevantPR != nil {
		e.RelevantPR = pkg.RelevantPR.URL()
	} else {
		err = errors.New("missing pr")
	}

	return
}

var NewSoftwarePkgAbandonedEvent = NewSoftwarePkgRejectedEvent

// softwarePkgAppliedEvent
type softwarePkgAppliedEvent struct {
	Importer          string `json:"importer"`
	ImporterEmail     string `json:"importer_email"`
	PkgId             string `json:"pkg_id"`
	PkgName           string `json:"pkg_name"`
	PkgDesc           string `json:"pkg_desc"`
	SourceCodeURL     string `json:"source_code_url"`
	SourceCodeLicense string `json:"source_code_license"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason_to_import"`
}

func (e *softwarePkgAppliedEvent) ToMessage() ([]byte, error) {
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
		SourceCodeURL:     app.SourceCode.Address.URL(),
		SourceCodeLicense: app.SourceCode.License.License(),
		ImportingPkgSig:   app.ImportingPkgSig.ImportingPkgSig(),
		ReasonToImportPkg: app.ReasonToImportPkg.ReasonToImportPkg(),
	}
}
