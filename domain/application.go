package domain

import "github.com/opensourceways/src-package-server/domain/dp"

type Application struct {
	SourceCodes     []dp.URL
	PackageDesc     dp.PackageDesc
	PackageLicense  dp.PackageLicense
	ReasonToImport  dp.ReasonToImportPkg
	PackagePlatform dp.PackagePlatform
	Sig             dp.Sig
}
