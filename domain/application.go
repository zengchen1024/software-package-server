package domain

import "github.com/opensourceways/src-package-server/domain/dp"

type Application struct {
	SourceCode        SourceCode
	PackageName       dp.PackageName
	PackageDesc       dp.PackageDesc
	PackagePlatform   dp.PackagePlatform
	ImportingPkgSig   dp.ImportingPkgSig
	ReasonToImportPkg dp.ReasonToImportPkg
}

type SourceCode struct {
	Address dp.URL
	License dp.License
}
