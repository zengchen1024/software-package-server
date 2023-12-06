package main

import (
	"time"

	kafka "github.com/opensourceways/kafka-lib/agent"
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/softwarepkgadapter"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/useradapterimpl"
	"github.com/opensourceways/software-package-server/watch/infrastructure/emailimpl"
	"github.com/opensourceways/software-package-server/watch/infrastructure/pullrequestimpl"
	watchrepoimpl "github.com/opensourceways/software-package-server/watch/infrastructure/repositoryimpl"
)

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type postgresqlConfig struct {
	DB         postgresql.Config    `json:"db"          required:"true"`
	Table      repositoryimpl.Table `json:"table"       required:"true"`
	WatchTable watchrepoimpl.Table  `json:"watch_table" require:"true"`
}

type watch struct {
	RobotToken     string `json:"robot_token"      required:"true"`
	CommunityOrg   string `json:"community_org"    required:"true"`
	CommunityRepo  string `json:"community_repo"   required:"true"`
	CISuccessLabel string `json:"ci_success_label" required:"true"`
	CIFailureLabel string `json:"ci_failure_label" required:"true"`
	// unit second
	Interval int `json:"interval"`
}

type mongoConfig struct {
	DB          mongdblib.Config               `json:"db"`
	Collections softwarepkgadapter.Collections `json:"collections"`
}

type Topics struct {
	SoftwarePkgInitialized string `json:"software_pkg_initialized" required:"true"`
}

// domainConfig
type domainConfig struct {
	domain.Config

	DomainPrimitive dp.Config `json:"domain_primitive"`
}

type Config struct {
	User        useradapterimpl.Config `json:"user"`
	Kafka       kafka.Config           `json:"kafka"`
	Watch       watch                  `json:"watch"`
	Email       emailimpl.Config       `json:"email"`
	Mongo       mongoConfig            `json:"mongo"`
	Topics      Topics                 `json:"topics"`
	Postgresql  postgresqlConfig       `json:"postgresql"`
	PullRequest pullrequestimpl.Config `json:"pull_request"`
	SoftwarePkg domainConfig           `json:"software_pkg"`
}

func loadConfig(path string) (*Config, error) {
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

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.User,
		&cfg.Kafka,
		&cfg.Watch,
		&cfg.PullRequest,
		&cfg.Mongo.DB,
		&cfg.Mongo.Collections,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Table,
		&cfg.SoftwarePkg.Config,
		&cfg.SoftwarePkg.DomainPrimitive,
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

func (w *watch) SetDefault() {
	if w.CommunityOrg == "" {
		w.CommunityOrg = "openeuler"
	}

	if w.CommunityRepo == "" {
		w.CommunityRepo = "community"
	}

	if w.CISuccessLabel == "" {
		w.CISuccessLabel = "ci_successful"
	}

	if w.CIFailureLabel == "" {
		w.CIFailureLabel = "ci_failed"
	}

	if w.Interval <= 0 {
		w.Interval = 10
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

func (w *watch) IntervalDuration() time.Duration {
	return time.Second * time.Duration(w.Interval)
}
