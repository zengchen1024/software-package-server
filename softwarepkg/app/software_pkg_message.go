package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/service"
	"github.com/sirupsen/logrus"
)

type SoftwarePkgMessageService interface {
	HandleCIChecking(cmd CmdToHandleCIChecking) error
}

type softwarePkgMessageService struct {
	repo         repository.SoftwarePkg
	message      message.SoftwarePkgMessage
	maintainer   maintainer.Maintainer
	reviewServie service.SoftwarePkgReviewService
}

func (s softwarePkgMessageService) HandleCIChecking(cmd CmdToHandleCIChecking) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	closedByFailedCI, err := pkg.HandleCI(cmd.isSuccess(), cmd.RelevantPR)
	if err != nil {
		return err
	}

	if pkg.Phase.IsClosed() {
		e := domain.NewSoftwarePkgAlreadyClosedEvent(cmd.RelevantPR)
		if err := s.message.NotifyPkgAlreadyClosed(&e); err != nil {
			logrus.Errorf(
				"failed to notify the pkg is already closed when handling ci checking, err:%s",
				err.Error(),
			)
		}

		if !closedByFailedCI {
			return nil
		}

		comment := s.genCommentOfFailedCI(&cmd)
		if err := s.repo.AddReviewComment(cmd.PkgId, &comment); err != nil {
			logrus.Errorf(
				"failed to add a comment when handling ci checking, err:%s",
				err.Error(),
			)
		}
	}

	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when handling ci checking, err:%s",
			err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) genCommentOfFailedCI(cmd *CmdToHandleCIChecking) domain.SoftwarePkgReviewComment {
	author, _ := dp.NewAccount("software-pkg-robot")

	str := fmt.Sprintf(
		"I'am sorry to close this application. Because the checking failed with the reason as bellow.\n\n%s",
		cmd.FiledReason,
	)
	content, _ := dp.NewReviewComment(str)

	return domain.NewSoftwarePkgReviewComment(author, content)
}
