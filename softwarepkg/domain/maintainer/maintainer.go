package maintainer

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type Maintainer interface {
	HasPermission(*domain.SoftwarePkg, *domain.User) (bool, bool)
	Reviewer(*domain.SoftwarePkg, *domain.User) domain.Reviewer
}
