package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

type SoftwarePkgReviewService interface {
	ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, user *domain.SoftwarePkgApprover) error
}

func NewReviewService(m message.SoftwarePkgMessage) SoftwarePkgReviewService {
	return &reviewService{message: m}
}

type reviewService struct {
	message message.SoftwarePkgMessage
}

func (s *reviewService) ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, user *domain.SoftwarePkgApprover) error {
	if approved, err := pkg.ApproveBy(user); !approved {
		return err
	}

	e := domain.NewSoftwarePkgApprovedEvent(pkg)
	err := s.message.NotifyPkgApproved(&e)
	s.log(pkg, "approved", err)

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
