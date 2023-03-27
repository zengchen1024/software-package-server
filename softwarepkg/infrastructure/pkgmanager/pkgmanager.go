package pkgmanager

import "github.com/opensourceways/robot-gitee-lib/client"

var instance *service

type service struct {
	cli client.Client
	org string
}

func Init(cfg *Config) {
	instance = &service{
		cli: client.NewClient(cfg.Token()),
		org: cfg.Org,
	}
}

func Instance() *service {
	return instance
}

func (s *service) IsPkgExisted(pkg string) bool {
	_, err := s.cli.GetRepo(s.org, pkg)

	return err == nil
}
