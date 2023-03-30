package maintainerimpl

import (
	"errors"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/cacheagent"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var instance *cacheagent.Agent

func Init(cfg *Config) error {
	v, err := cacheagent.NewCacheAgent(
		&sigLoader{
			cli:  utils.NewHttpClient(3),
			link: cfg.ReadURL,
		},
		cfg.IntervalDuration(),
	)

	if err == nil {
		instance = v
	}

	return err
}

func Exit() {
	if instance != nil {
		instance.Stop()
	}
}

func Maintainer() maintainerImpl {
	return maintainerImpl{instance}
}

// maintainerImpl
type maintainerImpl struct {
	agent *cacheagent.Agent
}

func (impl maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user *domain.User) bool {
	v := impl.agent.GetData()
	m, ok := v.(*sigData)

	return ok && m.hasMaintainer(user.GiteeID)
}

func (impl maintainerImpl) FindUser(giteeAccount string) (dp.Account, error) {
	return nil, errors.New("unimplemented")
}
