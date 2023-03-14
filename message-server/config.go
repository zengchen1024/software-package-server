package main

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
)

type Config struct {
	kafka.Config

	Topics    Topics `json:"topics"      required:"true"`
	GroupName string `json:"group_name"  required:"true"`
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	return cfg.Config.Validate()
}

type Topics struct {
	SoftwarePkgPRClosed    string `json:"software_pkg_pr_closed"      required:"true"`
	SoftwarePkgPRMerged    string `json:"software_pkg_pr_merged"      required:"true"`
	SoftwarePkgPRCIChecked string `json:"software_pkg_pr_ci_checked"  required:"true"`
	SoftwarePkgRepoCreated string `json:"software_pkg_repo_created"   required:"true"`
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
