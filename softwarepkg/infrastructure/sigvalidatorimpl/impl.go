package sigvalidatorimpl

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/cacheagent"
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

func SigValidator() sigValidatorImpl {
	return sigValidatorImpl{instance}
}

// sigValidatorImpl
type sigValidatorImpl struct {
	agent *cacheagent.Agent
}

func (impl sigValidatorImpl) IsValidSig(sig string) bool {
	v := impl.agent.GetData()
	m, ok := v.(*sigData)

	return ok && m.hasSig(sig)
}

func (impl sigValidatorImpl) GetAll() (info []sigDetail) {
	v := impl.agent.GetData()

	m, ok := v.(*sigData)
	if !ok {
		return nil
	}

	return m.getAll()
}
