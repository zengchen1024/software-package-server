package domain

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type SoftwarePkgApprovedEvent struct {
	PkgId       string
	Importer    dp.Account
	Application SoftwarePkgApplication
}

type SoftwarePkgRejectedEvent struct {
	PkgId    string
	Importer dp.Account
}
