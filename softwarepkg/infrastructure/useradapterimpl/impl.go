package useradapterimpl

import (
	"errors"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/cacheagent"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
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
}

func (impl *userAdapterImpl) HasPermission(info *domain.SoftwarePkg, user *domain.User) (
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

func (impl *userAdapterImpl) Find(giteeAccount string) (domain.User, error) {
	return domain.User{}, errors.New("unimplemented")
}

func (impl *userAdapterImpl) Roles(pkg *domain.SoftwarePkg, user *domain.User) (roles []dp.CommunityRole) {
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
