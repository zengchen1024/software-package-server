package app

import (
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/service"
	"github.com/sirupsen/logrus"
)

type SoftwarePkgService interface {
	ApplyNewPkg(*CmdToApplyNewSoftwarePkg) (string, error)
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
	message      message.SoftwarePkgMessage
	maintainer   maintainer.Maintainer
	reviewServie service.SoftwarePkgReviewService
}

func (s *softwarePkgService) ApplyNewPkg(cmd *CmdToApplyNewSoftwarePkg) (
	code string, err error,
) {
	v := domain.NewSoftwarePkg(cmd.Importer.Account, cmd.PkgName, &cmd.Application)
	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			code = errorSoftwarePkgExists
		}
	} else {
		e := domain.NewSoftwarePkgAppliedEvent(&cmd.Importer, &v)
		if err1 := s.message.NotifyPkgApplied(&e); err1 != nil {
			logrus.Errorf(
				"failed to notify a new applying pkg:%s, err:%s",
				v.Id, err.Error(),
			)
		} else {
			logrus.Debugf(
				"successfully to notify a new applying pkg:%s", v.Id,
			)
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
