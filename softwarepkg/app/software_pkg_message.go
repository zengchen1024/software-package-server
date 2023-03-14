package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgMessageService interface {
	HandleCIChecking(CmdToHandleCIChecking) error
	HandleRepoCreated(CmdToHandleRepoCreated) error
	HandlePkgRejected(CmdToHandlePkgRejected) error
}

type softwarePkgMessageService struct {
	repo       repository.SoftwarePkg
	message    message.SoftwarePkgMessage
	maintainer maintainer.Maintainer
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
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
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
			"failed to notify the pkg is already closed when %s, err:%s",
			cmd.logString(), err.Error(),
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
			"failed to add a comment when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}
}

// HandleRepoCreated
func (s softwarePkgMessageService) HandleRepoCreated(cmd CmdToHandleRepoCreated) error {
	if !cmd.isSuccess() {
		logrus.Errorf(
			"failed to create repo on platform:%s for pkg:%s, err:%s",
			cmd.Platform.PackagePlatform(), cmd.PkgId, cmd.FiledReason,
		)

		return nil
	}

	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	if err := pkg.HandleRepoCreated(cmd.RepoCreatedInfo); err != nil {
		return err
	}

	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) HandlePkgRejected(cmd CmdToHandlePkgRejected) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	if pkg.Phase.IsClosed() {
		return nil
	}

	user, err := s.maintainer.FindUser(cmd.RejectedBy)
	if err != nil {
		return err
	}

	has, err := s.maintainer.HasPermission(&pkg, user)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("no permission when %s", cmd.logString())
	}

	if _, err := pkg.RejectBy(user); err != nil {
		return err
	}

	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}
