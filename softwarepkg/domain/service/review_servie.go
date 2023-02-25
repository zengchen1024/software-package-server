package service

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/maintainer"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

type SoftwarePkgReviewService interface {
	ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account) error

	RejectPkg(pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account) error

	Close(pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account) error
}

type reviewService struct {
	maintainer maintainer.Maintainer
	repo       repository.SoftwarePkg
	message    message.SoftwarePkgMessage
}

func (s *reviewService) ApprovePkg(
	pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account,
) error {
	b, err := s.maintainer.HasPermission(pkg, user)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("no permission")
	}

	changed, approved := pkg.ApproveBy(user)
	if changed {
		return s.repo.SaveSoftwarePkg(pkg, version)
	}

	if approved {
		// TODO
		return s.message.NotifyPkgApproved(&domain.SoftwarePkgApprovedEvent{})
	}

	return nil
}

func (s *reviewService) RejectPkg(
	pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account,
) error {
	b, err := s.maintainer.HasPermission(pkg, user)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("no permission")
	}

	changed, rejected := pkg.RejectBy(user)
	if changed {
		return s.repo.SaveSoftwarePkg(pkg, version)
	}

	if rejected {
		// TODO
		return s.message.NotifyPkgRejected(&domain.SoftwarePkgRejectedEvent{})
	}

	return nil
}

func (s *reviewService) Close(
	pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account,
) error {
	b, err := s.maintainer.HasPermission(pkg, user)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("no permission")
	}

	if err := pkg.Close(); err != nil {
		return err
	}

	return s.repo.SaveSoftwarePkg(pkg, version)
}
