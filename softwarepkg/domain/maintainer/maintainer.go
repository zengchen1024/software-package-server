package maintainer

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Maintainer interface {
	HasPermission(*domain.SoftwarePkgBasicInfo, *domain.User) bool
	FindUser(string) (dp.Account, error)
}
