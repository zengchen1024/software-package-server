package sigvalidatorimpl

import "errors"

type Config struct {
	Sigs []string            `json:"sigs"`
	sigs map[string]struct{} `json:"-"`
}

func (cfg *Config) SetDefault() {}

func (cfg *Config) Validate() error {
	if len(cfg.Sigs) == 0 {
		return errors.New("empty sigs")
	}

	cfg.sigs = make(map[string]struct{}, len(cfg.Sigs))

	v := struct{}{}
	for _, item := range cfg.Sigs {
		cfg.sigs[item] = v
	}

	return nil
}

func (cfg *Config) hasSig(v string) bool {
	if cfg.sigs == nil {
		return false
	}

	_, ok := cfg.sigs[v]

	return ok
}
