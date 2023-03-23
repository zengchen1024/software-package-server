package sigvalidatorimpl

import "github.com/opensourceways/server-common-lib/config"

func Init(filepath string) (*sigValidatorImpl, error) {
	agent := config.NewConfigAgent(func() config.Config {
		return new(Config)
	})

	if err := agent.Start(filepath); err != nil {
		return nil, err
	}

	return &sigValidatorImpl{&agent}, nil
}

type sigValidatorImpl struct {
	agent *config.ConfigAgent
}

func (impl *sigValidatorImpl) IsValidSig(sig string) bool {
	_, v := impl.agent.GetConfig()

	cfg, ok := v.(*Config)

	return ok && cfg.hasSig(sig)
}

func (impl *sigValidatorImpl) Stop() {
	impl.agent.Stop()
}
