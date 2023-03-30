package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgci"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgMessageService interface {
	HandlePkgCIChecking(CmdToHandlePkgCIChecking) error
	HandlePkgCIChecked(CmdToHandlePkgCIChecked) error
	HandlePkgCodeSaved(CmdToHandlePkgCodeSaved) error
	HandlePkgInitialized(CmdToHandlePkgInitialized) error
	HandlePkgRepoCreated(CmdToHandlePkgRepoCreated) error
	HandlePkgAlreadyExisted(CmdToHandlePkgAlreadyExisted) error
}

func NewSoftwarePkgMessageService(
	ci pkgci.PkgCI,
	repo repository.SoftwarePkg,
	manager pkgmanager.PkgManager,
	message message.SoftwarePkgIndirectMessage,
) softwarePkgMessageService {
	return softwarePkgMessageService{
		ci:      ci,
		repo:    repo,
		manager: manager,
		message: message,
	}
}

type softwarePkgMessageService struct {
	ci      pkgci.PkgCI
	repo    repository.SoftwarePkg
	manager pkgmanager.PkgManager
	message message.SoftwarePkgIndirectMessage
}

// HandlePkgCIChecking
func (s softwarePkgMessageService) HandlePkgCIChecking(cmd CmdToHandlePkgCIChecking) error {
	return s.ci.SendTest(&cmd)
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

// HandlePkgCodeSaved
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

	if cmd.isPkgAreadyExisted() {
		if err := pkg.HandlePkgAlreadyExisted(); err != nil {
			return err
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

// HandlePkgAlreadyExisted
func (s softwarePkgMessageService) HandlePkgAlreadyExisted(cmd CmdToHandlePkgAlreadyExisted) error {
	if b, _ := s.repo.HasSoftwarePkg(cmd.PkgName); b {
		return nil
	}

	v, err := s.manager.GetPkg(cmd.PkgName)
	if err != nil {
		logrus.Errorf("get pkg/%s failed, err:%s", cmd.PkgName.PackageName(), err.Error())

		return err
	}

	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			return nil
		}

		logrus.Errorf(
			"failed to add a software pkg, pkgname:%s, err:%s",
			cmd.PkgName.PackageName(), err.Error(),
		)
	}

	return err
}
