package config

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
)

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
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

type PostgresqlConfig struct {
	DB postgresql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

type Config struct {
	MQ          messageimpl.Config `json:"mq"             required:"true"`
	Postgresql  PostgresqlConfig   `json:"postgresql"     required:"true"`
	SoftwarePkg dp.Config          `json:"software_pkg"   required:"true"`
	Api         middleware.Config  `json:"api"            required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.MQ,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.SoftwarePkg,
		&cfg.Api,
	}
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
