package pkgmanagerimpl

import (
	"github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var instance *service

func Init(cfg *Config) {
	instance = &service{
		cli: client.NewClient(cfg.Token()),
		org: cfg.Org,
	}
}

func Instance() *service {
	return instance
}

type service struct {
	cli client.Client
	org string
}

func (s *service) IsPkgExisted(pkg dp.PackageName) bool {
	_, err := s.cli.GetRepo(s.org, pkg.PackageName())

	return err == nil
}
