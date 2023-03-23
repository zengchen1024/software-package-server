package sigvalidatorimpl

import "errors"

type Config struct {
	Sigs []string `json:"sigs"`
}

func (cfg *Config) SetDefault() {}

func (cfg *Config) Validate() error {
	if len(cfg.Sigs) == 0 {
		return errors.New("empty sigs")
	}

	return nil
}

func (cfg *Config) hasSig(v string) bool {
	for _, item := range cfg.Sigs {
		if item == v {
			return true
		}
	}

	return false
}
