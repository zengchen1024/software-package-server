package app

import (
	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/service"
)

type softwarePkgMessageService struct {
	repo         repository.SoftwarePkg
	message      message.SoftwarePkgMessage
	maintainer   maintainer.Maintainer
	reviewServie service.SoftwarePkgReviewService
}

type CmdToHandleCIChecking struct {
	PkgId       string
	RelevantPR  dp.URL
	FiledReason string
}

func (s softwarePkgMessageService) HandleCIChecking(cmd CmdToHandleCIChecking) error {
	_, _, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	return nil
}
