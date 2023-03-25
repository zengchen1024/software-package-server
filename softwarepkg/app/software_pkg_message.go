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
	HandlePkgCIChecked(cmd CmdToHandlePkgCIChecked) error
	HandlePkgInitialized(cmd CmdToHandlePkgInitialized) error
	HandlePkgRepoCreated(CmdToHandlePkgRepoCreated) error
	HandlePkgCodeSaved(CmdToHandlePkgCodeSaved) error
}

func NewSoftwarePkgMessageService(
	repo repository.SoftwarePkg,
	message message.SoftwarePkgIndirectMessage,
	maintainer maintainer.Maintainer,
) softwarePkgMessageService {
	return softwarePkgMessageService{
		repo:       repo,
		message:    message,
		maintainer: maintainer,
	}
}

type softwarePkgMessageService struct {
	repo       repository.SoftwarePkg
	message    message.SoftwarePkgIndirectMessage
	maintainer maintainer.Maintainer
}

// HandlePkgCIChecked
func (s softwarePkgMessageService) HandlePkgCIChecked(cmd CmdToHandlePkgCIChecked) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	alreadyClosed, err := pkg.HandleCI(cmd.isSuccess(), cmd.RelevantPR)
	if err != nil || alreadyClosed {
		return err
	}

	if !cmd.isSuccess() {
		s.addCommentForFailedCI(&cmd)
	}

	pkg.PRNum = cmd.PRNum
	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) addCommentForFailedCI(cmd *CmdToHandlePkgCIChecked) {
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

// HandlePkgRepoCreated
func (s softwarePkgMessageService) HandlePkgRepoCreated(cmd CmdToHandlePkgRepoCreated) error {
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

// HandlePkgRepoCreated
func (s softwarePkgMessageService) HandlePkgCodeSaved(cmd CmdToHandlePkgCodeSaved) error {
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

	if err := pkg.HandleCodeSaved(); err != nil {
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

// HandlePkgInitDone
func (s softwarePkgMessageService) HandlePkgInitialized(cmd CmdToHandlePkgInitialized) error {
	pkg, version, err := s.repo.FindSoftwarePkgBasicInfo(cmd.PkgId)
	if err != nil {
		return err
	}

	if cmd.isSuccess() {
		if !pkg.Application.PackagePlatform.IsLocalPlatform() {
			s.notifyPkgInitialized(&pkg, &cmd)
		}

		return nil
	}

	if b := cmd.isPkgAreadyExisted(); b {
		if err := pkg.HandlePkgAlreadyExisted(); err != nil {
			return nil
		}

		s.addCommentForExistedPkg(&cmd)

		if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
			logrus.Errorf(
				"save pkg failed when %s, err:%s",
				cmd.logString(), err.Error(),
			)
		}

		return nil
	}

	logrus.Errorf("pkg init failed, pkgid:%s, err:%s", cmd.PkgId, cmd.FiledReason)

	return nil
}

func (s softwarePkgMessageService) notifyPkgInitialized(
	pkg *domain.SoftwarePkgBasicInfo, cmd *CmdToHandlePkgInitialized,
) {
	e := domain.NewSoftwarePkgInitializedEvent(pkg)

	if err := s.message.NotifyPkgIndirectlyApproved(&e); err != nil {
		logrus.Errorf(
			"failed to notify the pkg was approved indirectly when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}
}

func (s softwarePkgMessageService) addCommentForExistedPkg(cmd *CmdToHandlePkgInitialized) {
	author, _ := dp.NewAccount("software-pkg-robot")

	str := fmt.Sprintf(
		"I'am sorry to close this application. Because the pkg was imported sometimes ago. The repo address is %s. You can work on that repo.",
		cmd.RepoLink,
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

func (s softwarePkgMessageService) validateUser(pkg *domain.SoftwarePkgBasicInfo, v string) (
	dp.Account, error,
) {
	user, err := s.maintainer.FindUser(v)
	if err != nil {
		return nil, err
	}

	has, err := s.maintainer.HasPermission(pkg, user)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, fmt.Errorf("no permission")
	}

	return user, nil
}
