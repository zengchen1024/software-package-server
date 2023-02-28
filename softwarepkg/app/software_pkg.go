package app

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/service"
)

type SoftwarePkgService interface {
	ApplyNewPkg(dp.Account, *CmdToApplyNewSoftwarePkg) (string, error)
	GetPkgReviewDetail(string) (SoftwarePkgReviewDTO, error)
	ListPkgs(*CmdToListPkgs) (SoftwarePkgsDTO, error)
}

var _ SoftwarePkgService = (*softwarePkgService)(nil)

func NewSoftwarePkgService(repo repository.SoftwarePkg) *softwarePkgService {
	return &softwarePkgService{
		repo: repo,
	}
}

type softwarePkgService struct {
	repo         repository.SoftwarePkg
	maintainer   maintainer.Maintainer
	reviewServie service.SoftwarePkgReviewService
}

func (s *softwarePkgService) ApplyNewPkg(user dp.Account, cmd *CmdToApplyNewSoftwarePkg) (
	code string, err error,
) {
	v := domain.NewSoftwarePkg(user, (*domain.SoftwarePkgApplication)(cmd))
	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			code = errorSoftwarePkgExists
		}
	}

	return
}

func (s *softwarePkgService) ListPkgs(cmd *CmdToListPkgs) (SoftwarePkgsDTO, error) {
	v, total, err := s.repo.FindSoftwarePkgs(*cmd)
	if err != nil || len(v) == 0 {
		return SoftwarePkgsDTO{}, nil
	}

	return toSoftwarePkgsDTO(v, total), nil
}
