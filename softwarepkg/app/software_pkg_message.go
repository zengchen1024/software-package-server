package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgMessageService interface {
	HandleCIChecking(cmd CmdToHandleCIChecking) error
}

type softwarePkgMessageService struct {
	repo    repository.SoftwarePkg
	message message.SoftwarePkgMessage
}

func (s softwarePkgMessageService) HandleCIChecking(cmd CmdToHandleCIChecking) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	alreadyClosed, err := pkg.HandleCI(cmd.isSuccess(), cmd.RelevantPR)
	if err != nil {
		return err
	}

	if pkg.Phase.IsClosed() {
		if alreadyClosed {
			s.notifyPkgAlreadyClosed(&cmd)

			return nil
		}

		s.addCommentForFailedCI(&cmd)
	}

	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when handling ci checking, err:%s",
			err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) notifyPkgAlreadyClosed(cmd *CmdToHandleCIChecking) {
	if !cmd.isSuccess() {
		return
	}

	e := domain.NewSoftwarePkgAlreadyClosedEvent(cmd.PkgId, cmd.RelevantPR)
	if err := s.message.NotifyPkgAlreadyClosed(&e); err != nil {
		logrus.Errorf(
			"failed to notify the pkg is already closed when handling ci checking, err:%s",
			err.Error(),
		)
	}
}

func (s softwarePkgMessageService) addCommentForFailedCI(cmd *CmdToHandleCIChecking) {
	author, _ := dp.NewAccount("software-pkg-robot")

	str := fmt.Sprintf(
		"I'am sorry to close this application. Because the checking failed with the reason as bellow.\n\n%s",
		cmd.FiledReason,
	)
	content, _ := dp.NewReviewComment(str)

	comment := domain.NewSoftwarePkgReviewComment(author, content)

	if err := s.repo.AddReviewComment(cmd.PkgId, &comment); err != nil {
		logrus.Errorf(
			"failed to add a comment when handling ci checking, err:%s",
			err.Error(),
		)
	}
}
