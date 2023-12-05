package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgciimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/softwarepkgadapter"
	"github.com/opensourceways/software-package-server/utils"
)

func loadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.setDefault()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type domainConfig struct {
	domain.CIConfig

	DomainPrimitive dp.Config `json:"domain_primitive"`
}

type postgresqlConfig struct {
	DB    postgresql.Config    `json:"db"`
	Table repositoryimpl.Table `json:"table"`
}

type mongoConfig struct {
	DB          mongdblib.Config               `json:"db"`
	Collections softwarepkgadapter.Collections `json:"collections"`
}

type Config struct {
	CI          pkgciimpl.Config      `json:"ci"`
	Mongo       mongoConfig           `json:"mongo"`
	Kafka       kfklib.Config         `json:"kafka"`
	Topics      Topics                `json:"topics"`
	Postgresql  postgresqlConfig      `json:"postgresql"`
	Encryption  utils.Config          `json:"encryption"`
	PkgManager  pkgmanagerimpl.Config `json:"pkg_manager"`
	SoftwarePkg domainConfig          `json:"software_pkg"`
}

type Topics struct {
	SoftwarePkgCIDone         string `json:"software_pkg_ci_done"          required:"true"`
	SoftwarePkgApplied        string `json:"software_pkg_applied"          required:"true"`
	SoftwarePkgRetested       string `json:"software_pkg_retested"         required:"true"`
	SoftwarePkgRepoCodePushed string `json:"software_pkg_repo_code_pushed" required:"true"`
	SoftwarePkgAlreadyExisted string `json:"software_pkg_already_existed"  required:"true"`

	// importer edited the pkg and want to reload code
	SoftwarePkgCodeUpdated string `json:"software_pkg_code_updated"        required:"true"`

	// the pkg code has downloaded
	SoftwarePkgCodeChanged string `json:"software_pkg_code_changed"        required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.CI,
		&cfg.Mongo.DB,
		&cfg.Mongo.Collections,
		&cfg.Kafka,
		&cfg.Encryption,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Table,
		&cfg.PkgManager,
		&cfg.SoftwarePkg.CIConfig,
		&cfg.SoftwarePkg.DomainPrimitive,
	}
}

func (cfg *Config) setDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) validate() error {
	if err := utils.CheckConfig(cfg, ""); err != nil {
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
