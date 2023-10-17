package maintainerimpl

import (
	"errors"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/cacheagent"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var instance maintainerImpl

func Init(cfg *Config) error {
	v, err := cacheagent.NewCacheAgent(
		&sigLoader{
			cli:  utils.NewHttpClient(3),
			link: cfg.ReadURL,
		},
		cfg.IntervalDuration(),
	)

	if err != nil {
		return err
	}

	instance.agent = v
	instance.tcSig = cfg.TCSig

	return err
}

func Exit() {
	if instance.agent != nil {
		instance.agent.Stop()
	}
}

func Maintainer() *maintainerImpl {
	return &instance
}

// maintainerImpl
type maintainerImpl struct {
	agent *cacheagent.Agent
	tcSig string
}

func (impl *maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user *domain.User) (
	has bool, isTC bool,
) {
	v := impl.agent.GetData()
	m, ok := v.(*sigData)
	if !ok {
		return
	}

	if has = m.isSigMaintainer(user.GiteeID, impl.tcSig); has {
		isTC = true
	} else {
		has = m.isSigMaintainer(user.GiteeID, info.Sig())
	}

	return
}

func (impl *maintainerImpl) Reviewer(info *domain.SoftwarePkgBasicInfo, user *domain.User) domain.Reviewer {
	v := impl.agent.GetData()

	m, ok := v.(*sigData)
	if !ok {
		return domain.Reviewer{}
	}

	r := domain.Reviewer{
		User: user.Account,
	}

	if m.isSigMaintainer(user.GiteeID, impl.tcSig) {
		r.Role = append(r.Role, "tc")
	}

	if m.isSigMaintainer(user.GiteeID, info.Sig()) {
		r.Role = append(r.Role, "sig maintainer")
	}

	// TODO committer

	return r
}

func (impl *maintainerImpl) FindUser(giteeAccount string) (dp.Account, error) {
	return nil, errors.New("unimplemented")
}
