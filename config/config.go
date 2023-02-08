package config

import (
	"github.com/opensourceways/community-robot-lib/utils"
)

func LoadConfig(path string, cfg interface{}) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	if f, ok := cfg.(configSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(configValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type Config struct {
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{}
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
