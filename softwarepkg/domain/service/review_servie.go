package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

type SoftwarePkgReviewService interface {
	ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error
	RejectPkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error
	AbandonPkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error
}

type reviewService struct {
	message message.SoftwarePkgMessage
}

func (s *reviewService) ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error {
	approved, err := pkg.ApproveBy(user)
	if !approved {
		return err
	}

	op := "approved"
	if e, err := domain.NewSoftwarePkgApprovedEvent(pkg); err != nil {
		s.log(pkg, op, err)
	} else {
		err := s.message.NotifyPkgApproved(&e)
		s.log(pkg, op, err)
	}

	return nil
}

func (s *reviewService) RejectPkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error {
	rejected, err := pkg.RejectBy(user)
	if !rejected {
		return err
	}

	op := "rejected"
	if e, err := domain.NewSoftwarePkgRejectedEvent(pkg); err != nil {
		s.log(pkg, op, err)
	} else {
		err := s.message.NotifyPkgRejected(&e)
		s.log(pkg, op, err)
	}

	return nil
}

func (s *reviewService) AbandonPkg(pkg *domain.SoftwarePkgBasicInfo, user dp.Account) error {
	if err := pkg.Abandon(user); err != nil {
		return err
	}

	op := "abandoned"
	if e, err := domain.NewSoftwarePkgAbandonedEvent(pkg); err != nil {
		s.log(pkg, op, err)
	} else {
		err := s.message.NotifyPkgAbandoned(&e)
		s.log(pkg, op, err)
	}

	return nil
}

func (s *reviewService) log(pkg *domain.SoftwarePkgBasicInfo, op string, err error) {
	msg := fmt.Sprintf(
		"notify that a pkg:%s/%s was %s", pkg.Id, pkg.PkgName.PackageName(), op,
	)

	if err == nil {
		logrus.Debugf("successfully to %s", msg)
	} else {
		logrus.Errorf("failed to %s, err:%s", msg, err.Error())
	}
}
