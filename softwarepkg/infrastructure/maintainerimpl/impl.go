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

func (impl *maintainerImpl) HasPermission(info *domain.SoftwarePkg, user *domain.User) (
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
		has = m.isSigMaintainer(
			user.GiteeID, info.Sig.ImportingPkgSig(),
		)
	}

	return
}

func (impl *maintainerImpl) FindUser(giteeAccount string) (dp.Account, error) {
	return nil, errors.New("unimplemented")
}

func (impl *maintainerImpl) Roles(pkg *domain.SoftwarePkg, user *domain.User) (roles []dp.CommunityRole) {
	if pkg.IsCommitter(user) {
		roles = append(roles, dp.CommunityRoleCommitter, dp.CommunityRoleRepoMember)
	}

	v := impl.agent.GetData()

	m, ok := v.(*sigData)
	if !ok {
		return
	}

	if m.isSigMaintainer(user.GiteeID, impl.tcSig) {
		roles = append(roles, dp.CommunityRoleTC)
	}

	if m.isSigMaintainer(user.GiteeID, pkg.Sig.ImportingPkgSig()) {
		roles = append(roles, dp.CommunityRoleSigMaintainer, dp.CommunityRoleRepoMember)
	}

	return
}
