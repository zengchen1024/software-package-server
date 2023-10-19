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
	ApplyNewPkg(*CmdToApplyNewSoftwarePkg) (NewSoftwarePkgDTO, string, error)
	GetPkgReviewDetail(string) (SoftwarePkgReviewDTO, string, error)
	ListPkgs(*CmdToListPkgs) (SoftwarePkgsDTO, error)
	UpdateApplication(*CmdToUpdateSoftwarePkgApplication) (string, error)

	Review(pid string, user *domain.User, reviews []domain.CheckItemReviewInfo) (err error)
	Reject(string, *domain.User) error
	Abandon(string, *domain.User) (string, error)
	RerunCI(string, *domain.User) (string, error)
	NewReviewComment(string, *CmdToWriteSoftwarePkgReviewComment) (string, error)

	TranslateReviewComment(*CmdToTranslateReviewComment) (
		dto TranslatedReveiwCommentDTO, code string, err error,
	)
}

var (
	_ SoftwarePkgService = (*softwarePkgService)(nil)

	softwarePkgRobot = "software-pkg-robot"
)

func NewSoftwarePkgService(
	repo repository.SoftwarePkg,
	manager pkgmanager.PkgManager,
	message message.SoftwarePkgMessage,
	sensitive sensitivewords.SensitiveWords,
	maintainer maintainer.Maintainer,
	translation translation.Translation,
	commentRepo repository.SoftwarePkgComment,
) *softwarePkgService {
	robot, _ := dp.NewAccount(softwarePkgRobot)

	return &softwarePkgService{
		repo:        repo,
		robot:       robot,
		message:     message,
		sensitive:   sensitive,
		maintainer:  maintainer,
		translation: translation,
		pkgService:  service.NewPkgService(manager, message),
		commentRepo: commentRepo,
	}
}

type softwarePkgService struct {
	repo        repository.SoftwarePkg
	robot       dp.Account
	message     message.SoftwarePkgMessage
	sensitive   sensitivewords.SensitiveWords
	pkgService  service.SoftwarePkgService
	maintainer  maintainer.Maintainer
	translation translation.Translation
	commentRepo repository.SoftwarePkgComment
}

func (s *softwarePkgService) ApplyNewPkg(cmd *CmdToApplyNewSoftwarePkg) (
	dto NewSoftwarePkgDTO, code string, err error,
) {
	v := domain.NewSoftwarePkg(
		cmd.Sig, &cmd.Repo, &cmd.Code, &cmd.Basic, &cmd.Importer,
	)
	if s.pkgService.IsPkgExisted(cmd.Basic.Name) {
		err = errors.New("software package already existed")
		code = errorSoftwarePkgExists

		return
	}

	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			code = errorSoftwarePkgExists
		}
	} else {
		dto.Id = v.Id

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

func (s *softwarePkgService) UpdateApplication(cmd *CmdToUpdateSoftwarePkgApplication) (string, error) {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return errorCodeForFindingPkg(err), err
	}

	err = pkg.UpdateApplication(&cmd.Basic, cmd.Sig, &cmd.Repo, &cmd.Importer)
	if err != nil {
		return domain.ParseErrorCode(err), err
	}

	err = s.repo.SaveSoftwarePkg(&pkg, version)

	return "", err
}
