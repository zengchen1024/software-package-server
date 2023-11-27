package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/sirupsen/logrus"
)

type SoftwarePkgInitAppService interface {
	ListApprovedPkgs() ([]string, error)
	SoftwarePkg(pkgId string) (domain.SoftwarePkg, error)
	HandlePkgInitDone(pkgId string, pr dp.URL) error
	HandlePkgInitStarted(pkgId string, pr dp.URL) error
	HandlePkgAlreadyExisted(pkgId string, repoLink string) error
}

var _ SoftwarePkgInitAppService = (*softwarePkgInitAppService)(nil)

func NewSoftwarePkgInitAppService(
	repo repository.SoftwarePkg,
	manager pkgmanager.PkgManager,
	message message.SoftwarePkgInitMessage,
	commentRepo repository.SoftwarePkgComment,
) *softwarePkgInitAppService {
	robot, _ := dp.NewAccount(softwarePkgRobot)

	return &softwarePkgInitAppService{
		repo:        repo,
		robot:       robot,
		message:     message,
		commentRepo: commentRepo,
	}
}

type softwarePkgInitAppService struct {
	repo        repository.SoftwarePkg
	robot       dp.Account
	message     message.SoftwarePkgInitMessage
	commentRepo repository.SoftwarePkgComment
}

// return pkg ids
func (s *softwarePkgInitAppService) ListApprovedPkgs() ([]string, error) {
	return s.repo.FindAllApproved(dp.PackagePhaseCreatingRepo)
}

func (s *softwarePkgInitAppService) SoftwarePkg(pkgId string) (domain.SoftwarePkg, error) {
	pkg, _, err := s.repo.Find(pkgId)

	return pkg, err
}

func (s *softwarePkgInitAppService) HandlePkgInitStarted(pkgId string, pr dp.URL) error {
	pkg, version, err := s.repo.FindAndIgnoreReview(pkgId)
	if err != nil {
		return err
	}

	if err := pkg.StartInitialization(pr); err != nil {
		return err
	}

	return s.repo.SaveAndIgnoreReview(&pkg, version)
}

func (s *softwarePkgInitAppService) HandlePkgInitDone(pkgId string, pr dp.URL) error {
	pkg, version, err := s.repo.FindAndIgnoreReview(pkgId)
	if err != nil {
		return err
	}

	if err := pkg.HandleInitialized(pr); err != nil {
		return err
	}

	if err := s.repo.SaveAndIgnoreReview(&pkg, version); err != nil {
		return err
	}

	s.notifyPkgInitialized(&pkg)

	return nil
}

func (s *softwarePkgInitAppService) HandlePkgAlreadyExisted(pkgId, repoLink string) error {
	pkg, version, err := s.repo.FindAndIgnoreReview(pkgId)
	if err != nil {
		return err
	}

	if err := pkg.HandleAlreadyExisted(); err != nil {
		return err
	}

	if err := s.repo.SaveAndIgnoreReview(&pkg, version); err != nil {
		return err
	}

	s.addCommentForExistedPkg(&pkg, repoLink)

	return nil
}

func (s *softwarePkgInitAppService) addCommentForExistedPkg(pkg *domain.SoftwarePkg, repoLink string) {
	str := fmt.Sprintf(
		"I'am sorry to close this application. Because the pkg was imported sometimes ago. The repo address is %s. You can work on that repo.",
		repoLink,
	)
	content, _ := dp.NewReviewComment(str)
	comment := domain.NewSoftwarePkgReviewComment(s.robot, content)

	if err := s.commentRepo.AddReviewComment(pkg.Id, &comment); err != nil {
		logrus.Errorf(
			"failed to add a comment when handling pkg already existed, pkg:%s, err:%s",
			pkg.Id, err.Error(),
		)
	}
}

func (s *softwarePkgInitAppService) notifyPkgInitialized(pkg *domain.SoftwarePkg) {
	e := domain.NewSoftwarePkgInitializedEvent(pkg)

	if err := s.message.SendPkgInitializedEvent(&e); err != nil {
		logrus.Errorf(
			"failed to notify that the pkg was initialized, pkg: %s, err:%s",
			pkg.Id, err.Error(),
		)
	}
}
