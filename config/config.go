package config

import (
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"

	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/clavalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sensitivewordsimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/softwarepkgadapter"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/translationimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/useradapterimpl"
	"github.com/opensourceways/software-package-server/utils"
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

type domainConfig struct {
	domain.Config

	DomainPrimitive dp.Config `json:"domain_primitive"  required:"true"`
}

type postgresqlConfig struct {
	DB postgresql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

type mongoConfig struct {
	DB          mongdblib.Config          `json:"db"`
	Collections softwarepkgadapter.Config `json:"collections"`
}

type Config struct {
	MQ             messageimpl.Config        `json:"mq"                   required:"true"`
	CLA            clavalidatorimpl.Config   `json:"cla"                  required:"true"`
	User           useradapterimpl.Config    `json:"user"                 required:"true"`
	Mongo          mongoConfig               `json:"mongo"                 required:"true"`
	Encryption     utils.Config              `json:"encryption"           required:"true"`
	PkgManager     pkgmanagerimpl.Config     `json:"pkg_manager"          required:"true"`
	Middleware     middleware.Config         `json:"middleware"           required:"true"`
	Postgresql     postgresqlConfig          `json:"postgresql"           required:"true"`
	SoftwarePkg    domainConfig              `json:"software_pkg"         required:"true"`
	Translation    translationimpl.Config    `json:"translation"          required:"true"`
	SigValidator   sigvalidatorimpl.Config   `json:"sig"                  required:"true"`
	SensitiveWords sensitivewordsimpl.Config `json:"sensitive_words"      required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.MQ,
		&cfg.CLA,
		&cfg.User,
		&cfg.Mongo,
		&cfg.Mongo.Collections,
		&cfg.Encryption,
		&cfg.PkgManager,
		&cfg.Middleware,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.SoftwarePkg.Config,
		&cfg.SoftwarePkg.DomainPrimitive,
		&cfg.Translation,
		&cfg.SigValidator,
		&cfg.SensitiveWords,
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
