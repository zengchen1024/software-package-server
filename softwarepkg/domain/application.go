package domain

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageName       dp.PackageName
	PackageDesc       dp.PackageDesc
	PackagePlatform   dp.PackagePlatform
	ImportingPkgSig   dp.ImportingPkgSig
	ReasonToImportPkg dp.ReasonToImportPkg
}

type SoftwarePkgSourceCode struct {
	Address dp.URL
	License dp.License
}
