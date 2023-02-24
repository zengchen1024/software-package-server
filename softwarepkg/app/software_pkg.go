package app

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgService interface {
	ApplyNewSoftwarePkg(dp.Account, *CmdToApplyNewSoftwarePkg) (string, error)
}

func NewSoftwarePkgService(repo repository.SoftwarePkg) *softwarePkgService {
	return &softwarePkgService{repo}
}

type softwarePkgService struct {
	repo repository.SoftwarePkg
}

func (s *softwarePkgService) ApplyNewSoftwarePkg(
	user dp.Account, cmd *CmdToApplyNewSoftwarePkg,
) (code string, err error) {
	v := domain.NewSoftwarePkg(user, (*domain.Application)(cmd))

	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			code = errorSoftwarePkgExists
		}
	}

	return
}
