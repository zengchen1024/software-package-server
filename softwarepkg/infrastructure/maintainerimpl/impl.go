package maintainerimpl

import (
	"errors"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var instance *agent

func Init(cfg *Config) error {
	v := &agent{
		t: utils.NewTimer(),
		loader: sigLoader{
			cli: utils.NewHttpClient(3),
		},
	}

	err := v.start(cfg.ReadURL, cfg.IntervalDuration())
	if err == nil {
		instance = v
	}

	return err
}

func Exit() {
	if instance != nil {
		instance.stop()
	}
}

func Maintainer() maintainerImpl {
	return maintainerImpl{instance}
}

// maintainerImpl
type maintainerImpl struct {
	agent *agent
}

func (impl maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user *domain.User) bool {
	v := impl.agent.getSigData()

	return v.hasMaintainer(user.GiteeID)
}

func (impl maintainerImpl) FindUser(giteeAccount string) (dp.Account, error) {
	return nil, errors.New("unimplemented")
}
