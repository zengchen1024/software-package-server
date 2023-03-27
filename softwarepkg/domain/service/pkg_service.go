package service

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
	pkgmanager "github.com/opensourceways/software-package-server/softwarepkg/domain/pkgmanager"
)

type SoftwarePkgService interface {
	IsPkgExisted(string) bool
}

func NewPkgService(cli pkgmanager.PkgManager, message message.SoftwarePkgMessage) SoftwarePkgService {
	return &pkgService{cli: cli, message: message}
}

type pkgService struct {
	cli     pkgmanager.PkgManager
	message message.SoftwarePkgMessage
}

func (p *pkgService) IsPkgExisted(pkg string) bool {
	if !p.cli.IsPkgExisted(pkg) {
		return false
	}

	e := domain.NewSoftwarePkgAlreadyExistEvent(pkg)

	if err := p.message.NotifyPkgAlreadyExisted(&e); err != nil {
		logrus.Errorf(
			"failed to notify pkg already exist event,err:%s", err.Error(),
		)
	} else {
		logrus.Debugf("successfully to notify pkg already exist event")
	}

	return true
}
