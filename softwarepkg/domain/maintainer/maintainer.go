package maintainer

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Maintainer interface {
	Roles(*domain.SoftwarePkg, *domain.User) []dp.CommunityRole
}
