package domain

import (
	"encoding/json"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type SoftwarePkgApprovedEvent struct {
	PkgId       string
	Importer    dp.Account
	Application SoftwarePkgApplication
}

func (e *SoftwarePkgApprovedEvent) ToMessage() ([]byte, error) {
	return nil, nil
}

type SoftwarePkgRejectedEvent struct {
	PkgId    string
	Importer dp.Account
}

func (e *SoftwarePkgRejectedEvent) ToMessage() ([]byte, error) {
	return nil, nil
}

type SoftwarePkgAppliedEvent struct {
	Importer      string `json:"importer"`
	ImporterEmail string `json:"importer_email"`
	PkgId         string `json:"pkg_id"`
	PkgName       string `json:"pkg_name"`
	PkgDesc       string `json:"pkg_desc"`
	SourceCodeURL string `json:"source_code_url"`
}

func (e *SoftwarePkgAppliedEvent) ToMessage() ([]byte, error) {
	return json.Marshal(e)
}

func NewSoftwarePkgAppliedEvent(
	importer *User,
	pkg *SoftwarePkgBasicInfo,
) SoftwarePkgAppliedEvent {
	app := &pkg.Application

	return SoftwarePkgAppliedEvent{
		Importer:      importer.Account.Account(),
		ImporterEmail: importer.Email.Email(),
		PkgId:         pkg.Id,
		PkgName:       pkg.PkgName.PackageName(),
		PkgDesc:       app.PackageDesc.PackageDesc(),
		SourceCodeURL: app.SourceCode.Address.URL(),
	}
}
