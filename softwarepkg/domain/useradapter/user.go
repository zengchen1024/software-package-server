package useradapter

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	//"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type UserAdapter interface {
	Find(pid string, platform string) (domain.User, error)

	//Roles(*domain.SoftwarePkg, *domain.Reviewer) []dp.CommunityRole
}
