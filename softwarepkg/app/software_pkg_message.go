package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgcodeadapter"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgMessageService interface {
	DownloadPkgCode(cmd CmdToDownloadPkgCode) error
	StartCI(cmd CmdToStartCI) error
	HandlePkgCIDone(CmdToHandlePkgCIDone) error
	HandlePkgRepoCodePushed(CmdToHandlePkgRepoCodePushed) error
	//HandlePkgInitialized(CmdToHandlePkgInitialized) error
	ImportPkg(CmdToHandlePkgAlreadyExisted) error
}

func NewSoftwarePkgMessageService(
	code pkgcodeadapter.PkgCodeAdapter,
	repo repository.SoftwarePkg,
	manager pkgmanager.PkgManager,
	message message.SoftwarePkgIndirectMessage,
	commentRepo repository.SoftwarePkgComment,
) softwarePkgMessageService {
	robot, _ := dp.NewAccount(softwarePkgRobot)

	return softwarePkgMessageService{
		code:        code,
		repo:        repo,
		robot:       robot,
		manager:     manager,
		message:     message,
		commentRepo: commentRepo,
	}
}

type softwarePkgMessageService struct {
	code        pkgcodeadapter.PkgCodeAdapter
	repo        repository.SoftwarePkg
	robot       dp.Account
	manager     pkgmanager.PkgManager
	message     message.SoftwarePkgIndirectMessage
	commentRepo repository.SoftwarePkgComment
}

// DownloadPkgCode
func (s softwarePkgMessageService) DownloadPkgCode(cmd CmdToDownloadPkgCode) error {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	files, err := s.code.Download(&pkg)
	if err != nil {
		return err
	}

	pkg1, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	if !pkg1.SaveDownloadedFiles(files) {
		return nil
	}

	if err = s.repo.SaveSoftwarePkg(&pkg1, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)

		return err
	}

	s.notifyPkgCodeChanged(&pkg1)

	return nil
}

func (s softwarePkgMessageService) notifyPkgCodeChanged(pkg *domain.SoftwarePkg) {
	e := domain.NewSoftwarePkgCodeChangeedEvent(pkg)

	if err := s.message.SendSoftwarePkgCodeChangedEvent(&e); err != nil {
		logrus.Errorf(
			"failed to send pkg code changed event, pkg:%s, err:%s",
			pkg.Id, err.Error(),
		)
	}
}

// StartCI
func (s softwarePkgMessageService) StartCI(cmd CmdToStartCI) error {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	if err = pkg.StartCI(); err != nil {
		return err
	}

	if err = s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}

// HandlePkgCIDone
func (s softwarePkgMessageService) HandlePkgCIDone(cmd CmdToHandlePkgCIDone) error {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	if err := pkg.HandleCIDone(cmd.PRNumber, cmd.Success); err != nil {
		return err
	}

	s.addCIComment(&cmd)

	if err = s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) addCIComment(cmd *CmdToHandlePkgCIDone) {
	content, _ := dp.NewReviewComment(cmd.Detail)
	comment := domain.NewSoftwarePkgReviewComment(s.robot, content)

	if err := s.commentRepo.AddReviewComment(cmd.PkgId, &comment); err != nil {
		logrus.Errorf(
			"failed to add a comment when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}
}

// HandlePkgRepoCodePushed
func (s softwarePkgMessageService) HandlePkgRepoCodePushed(cmd CmdToHandlePkgRepoCodePushed) error {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	if err := pkg.HandleRepoCodePushed(); err != nil {
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

// HandlePkgInitialized
func (s softwarePkgMessageService) HandlePkgInitialized(cmd CmdToHandlePkgInitialized) error {
	pkg, version, err := s.repo.FindSoftwarePkg(cmd.PkgId)
	if err != nil {
		return err
	}

	if cmd.isSuccess() {
		if err := pkg.HandleInitialized(cmd.RelevantPR); err != nil {
			return err
		}

		if !pkg.Repo.Platform.IsLocalPlatform() {
			//s.notifyPkgInitialized(&pkg, &cmd)
		}
	} else {
		if !cmd.isPkgAreadyExisted() {
			logrus.Errorf("pkg init failed, pkgid:%s, err:%s", cmd.PkgId, cmd.FiledReason)

			return nil
		}

		if err := pkg.HandleAlreadyExisted(); err != nil {
			return err
		}

		s.addCommentForExistedPkg(&cmd)
	}

	if err := s.repo.SaveSoftwarePkg(&pkg, version); err != nil {
		logrus.Errorf(
			"save pkg failed when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return nil
}

func (s softwarePkgMessageService) addCommentForExistedPkg(cmd *CmdToHandlePkgInitialized) {
	str := fmt.Sprintf(
		"I'am sorry to close this application. Because the pkg was imported sometimes ago. The repo address is %s. You can work on that repo.",
		cmd.RepoLink,
	)
	content, _ := dp.NewReviewComment(str)
	comment := domain.NewSoftwarePkgReviewComment(s.robot, content)

	if err := s.commentRepo.AddReviewComment(cmd.PkgId, &comment); err != nil {
		logrus.Errorf(
			"failed to add a comment when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}
}

// ImportPkg
func (s softwarePkgMessageService) ImportPkg(cmd CmdToHandlePkgAlreadyExisted) error {
	v, err := s.manager.GetPkg(cmd.PkgName)
	if err != nil {
		logrus.Errorf(
			"failed to get pkg info when %s, err:%s",
			cmd.logString(), err.Error(),
		)

		return err
	}

	if err = s.repo.AddSoftwarePkg(&v); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			return nil
		}

		logrus.Errorf(
			"failed to add a software pkg when %s, err:%s",
			cmd.logString(), err.Error(),
		)
	}

	return err
}
