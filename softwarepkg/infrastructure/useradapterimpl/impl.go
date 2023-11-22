package useradapterimpl

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/cacheagent"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

var instance userAdapterImpl

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
	instance.om = omClient{cfg.OM}

	return err
}

func Exit() {
	if instance.agent != nil {
		instance.agent.Stop()
	}
}

func UserAdapter() *userAdapterImpl {
	return &instance
}

// userAdapterImpl
type userAdapterImpl struct {
	agent *cacheagent.Agent
	tcSig string

	om omClient
}

func (impl *userAdapterImpl) Find(pid, platform string) (domain.User, error) {
	return impl.om.getUserInfo(pid, platform)
}

func (impl *userAdapterImpl) Roles(pkg *domain.SoftwarePkg, user *domain.Reviewer) (tc bool, sigMaitainer bool) {
	if user.GiteeID == "" {
		return
	}

	v := impl.agent.GetData()

	m, ok := v.(*sigData)
	if !ok {
		return
	}

	tc = m.isSigMaintainer(user.GiteeID, impl.tcSig)

	sigMaitainer = m.isSigMaintainer(user.GiteeID, pkg.Sig.ImportingPkgSig())

	return
}
