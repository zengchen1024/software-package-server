package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgciimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	localutils "github.com/opensourceways/software-package-server/utils"
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

type domainConfig struct {
	domain.Config

	DomainPrimitive dp.Config `json:"domain_primitive"  required:"true"`
}

type postgresqlConfig struct {
	DB postgresql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

type Config struct {
	Kafka        kfklib.Config           `json:"kafka"                required:"true"`
	Topics       Topics                  `json:"topics_to_subscribe"  required:"true"`
	GroupName    string                  `json:"group_name"           required:"true"`
	Encryption   localutils.Config       `json:"encryption"           required:"true"`
	Postgresql   postgresqlConfig        `json:"postgresql"           required:"true"`
	PkgManager   pkgmanagerimpl.Config   `json:"pkg_manager"          required:"true"`
	SoftwarePkg  domainConfig            `json:"software_pkg"         required:"true"`
	SigValidator sigvalidatorimpl.Config `json:"sig"                  required:"true"`
	PkgCI        pkgciimpl.Config        `json:"ci"                   required:"true"`
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

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Kafka,
		&cfg.Encryption,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.SoftwarePkg.Config,
		&cfg.SoftwarePkg.DomainPrimitive,
		&cfg.PkgManager,
		&cfg.SigValidator,
		&cfg.PkgCI,
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
