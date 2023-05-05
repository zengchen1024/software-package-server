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
	instance.cfg = cfg.ConfigForPermission

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
	cfg   ConfigForPermission
}

func (impl *maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user *domain.User) (has bool, isTC bool) {
	v := impl.agent.GetData()
	m, ok := v.(*sigData)
	if !ok {
		return
	}

	if has = m.isSigMaintainer(user.GiteeID, impl.cfg.TCSig); has {
		isTC = true
	} else {
		has = m.isSigMaintainer(user.GiteeID,
			info.Application.ImportingPkgSig.ImportingPkgSig(),
		)
	}

	return
}

func (impl *maintainerImpl) HasPassedReview(info *domain.SoftwarePkgBasicInfo) bool {
	sig := info.Application.ImportingPkgSig.ImportingPkgSig()
	if sig == impl.cfg.EcoPkgSig {
		return len(info.ApprovedBy) > 0
	}

	numApprovedByTc := 0
	numApprovedBySigMaintainer := 0

	for i := range info.ApprovedBy {
		if info.ApprovedBy[i].IsTC {
			numApprovedByTc++
			numApprovedBySigMaintainer++
		} else {
			numApprovedBySigMaintainer++
		}
	}

	c := numApprovedByTc >= impl.cfg.MinNumApprovedByTC
	c1 := numApprovedBySigMaintainer >= impl.cfg.MinNumApprovedBySigMaintainer

	return c && c1
}

func (impl *maintainerImpl) FindUser(giteeAccount string) (dp.Account, error) {
	return nil, errors.New("unimplemented")
}
