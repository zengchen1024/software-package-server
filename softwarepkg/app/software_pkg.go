package app

import (
	"errors"

	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/sensitivewords"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/service"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/translation"
)

type SoftwarePkgService interface {
	ApplyNewPkg(*CmdToApplyNewSoftwarePkg) (string, error)
	GetPkgReviewDetail(string) (SoftwarePkgReviewDTO, error)
	ListPkgs(*CmdToListPkgs) (SoftwarePkgsDTO, error)
	UpdateApplication(*CmdToUpdateSoftwarePkgApplication) error

	Approve(string, dp.Account) (string, error)
	Reject(string, dp.Account) (string, error)
	Abandon(string, dp.Account) (string, error)
	NewReviewComment(string, *CmdToWriteSoftwarePkgReviewComment) (string, error)

	TranslateReviewComment(*CmdToTranslateReviewComment) (
		dto TranslatedReveiwCommentDTO, code string, err error,
	)
}

var _ SoftwarePkgService = (*softwarePkgService)(nil)

func NewSoftwarePkgService(
	repo repository.SoftwarePkg,
	manager pkgmanager.PkgManager,
	message message.SoftwarePkgMessage,
	sensitive sensitivewords.SensitiveWords,
	maintainer maintainer.Maintainer,
	translation translation.Translation,
) *softwarePkgService {
	return &softwarePkgService{
		repo:         repo,
		message:      message,
		sensitive:    sensitive,
		maintainer:   maintainer,
		translation:  translation,
		pkgService:   service.NewPkgService(manager, message),
		reviewServie: service.NewReviewService(message),
	}
}

type softwarePkgService struct {
	repo         repository.SoftwarePkg
	message      message.SoftwarePkgMessage
	sensitive    sensitivewords.SensitiveWords
	maintainer   maintainer.Maintainer
	translation  translation.Translation
	pkgService   service.SoftwarePkgService
	reviewServie service.SoftwarePkgReviewService
}

func (s *softwarePkgService) ApplyNewPkg(cmd *CmdToApplyNewSoftwarePkg) (
	code string, err error,
) {
	v := domain.NewSoftwarePkg(&cmd.Importer, cmd.PkgName, &cmd.Application)
	if s.pkgService.IsPkgExisted(cmd.PkgName) {
		err = errors.New("software package already existed")
		code = errorSoftwarePkgExists

		return
	}

	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			code = errorSoftwarePkgExists
		}
	} else {
		e := domain.NewSoftwarePkgAppliedEvent(&v)
		if err1 := s.message.NotifyPkgApplied(&e); err1 != nil {
			logrus.Errorf(
				"failed to notify a new applying pkg:%s, err:%s",
				v.Id, err1.Error(),
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

func (s *softwarePkgService) UpdateApplication(cmd *CmdToUpdateSoftwarePkgApplication) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	if err = pkg.UpdateApplication(&cmd.Application, &cmd.Importer); err != nil {
		return err
	}

	return s.repo.SaveSoftwarePkg(&pkg, version)
}
