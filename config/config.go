package config

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"

	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/controller"
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

// domainConfig
type domainConfig struct {
	domain.Config

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

type kafkaConfig struct {
	kfklib.Config

	Topics messageimpl.Topics `json:"topics"`
}

type Config struct {
	MQ             kafkaConfig               `json:"mq"`
	API            controller.Config         `json:"api"`
	CLA            clavalidatorimpl.Config   `json:"cla"`
	User           useradapterimpl.Config    `json:"user"`
	Mongo          mongoConfig               `json:"mongo"`
	Encryption     utils.Config              `json:"encryption"`
	PkgManager     pkgmanagerimpl.Config     `json:"pkg_manager"`
	Middleware     middleware.Config         `json:"middleware"`
	Postgresql     postgresqlConfig          `json:"postgresql"`
	SoftwarePkg    domainConfig              `json:"software_pkg"`
	Translation    translationimpl.Config    `json:"translation"`
	SigValidator   sigvalidatorimpl.Config   `json:"sig"`
	SensitiveWords sensitivewordsimpl.Config `json:"sensitive_words"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.MQ,
		&cfg.API,
		&cfg.CLA,
		&cfg.User,
		&cfg.Mongo,
		&cfg.Mongo.Collections,
		&cfg.Encryption,
		&cfg.PkgManager,
		&cfg.Middleware,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Table,
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
