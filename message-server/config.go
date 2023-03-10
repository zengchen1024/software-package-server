package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
)

type Config struct {
	kafka.Config

	Topics Topics `json:"topics"  required:"true"`
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	return cfg.Config.Validate()
}

type Topics struct {
	SoftwarePkgCIChecking string `json:"software_pkg_ci_checking"   required:"true"`
}

func loadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
