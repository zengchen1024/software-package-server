package maintainer

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Maintainer interface {
	HasPermission(*domain.SoftwarePkgBasicInfo, dp.Account) (bool, error)
}
