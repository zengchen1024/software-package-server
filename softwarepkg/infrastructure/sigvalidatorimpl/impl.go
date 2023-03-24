package sigvalidatorimpl

import "github.com/opensourceways/server-common-lib/utils"

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
		instance.Stop()
	}
}

func SigValidator() sigValidatorImpl {
	return sigValidatorImpl{instance}
}

// sigValidatorImpl
type sigValidatorImpl struct {
	agent *agent
}

func (impl sigValidatorImpl) IsValidSig(sig string) bool {
	v := impl.agent.getSigData()

	return v.hasSig(sig)
}

func (impl sigValidatorImpl) GetAll() (info []sigDetail) {
	v := impl.agent.getSigData()

	return v.getAll()
}
