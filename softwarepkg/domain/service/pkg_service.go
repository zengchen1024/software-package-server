package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	pkgmanager "github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
)

type SoftwarePkgService interface {
	IsPkgExisted(dp.PackageName) bool
}

func NewPkgService(
	manager pkgmanager.PkgManager, message message.SoftwarePkgMessage,
) SoftwarePkgService {
	return &pkgService{
		manager: manager,
		message: message,
	}
}

type pkgService struct {
	manager pkgmanager.PkgManager
	message message.SoftwarePkgMessage
}

func (s *pkgService) IsPkgExisted(pkg dp.PackageName) bool {
	if !s.manager.IsPkgExisted(pkg) {
		return false
	}

	e := domain.NewSoftwarePkgAlreadyExistEvent(pkg)
	err := s.message.NotifyPkgAlreadyExisted(&e)
	s.log(pkg, err)

	return true
}

func (s *pkgService) log(pkg dp.PackageName, err error) {
	msg := fmt.Sprintf(
		"notify that a pkg:%s already existed", pkg.PackageName(),
	)

	if err == nil {
		logrus.Debugf("successfully to %s", msg)
	} else {
		logrus.Errorf("failed to %s, err:%s", msg, err.Error())
	}
}
