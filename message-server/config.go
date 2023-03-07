package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
)

type Config struct {
	kafka.Config

	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	SoftwarePkgCIPassed string `json:"software_pkg_ci_passed"   required:"true"`
}

func loadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
