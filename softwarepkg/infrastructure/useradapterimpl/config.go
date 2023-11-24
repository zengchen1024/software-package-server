package useradapterimpl

import (
	"time"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Config struct {
	TCSig   string `json:"tc_sig"`
	ReadURL string `json:"read_url" required:"true"`

	// Interval the unit is hour
	Interval int `json:"interval" required:"true"`

	OM omConfig `json:"om"`
}

func (cfg *Config) SetDefault() {
	if cfg.TCSig == "" {
		cfg.TCSig = "TC"
	}
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Duration(cfg.Interval) * time.Hour
}

type omConfig struct {
	AppId              string `json:"app_id"                 required:"true"`
	AppSecret          string `json:"app_secret"             required:"true"`
	TokenEndpoint      string `json:"token_endpoint"         required:"true"`
	GiteeUserEndpoint  string `json:"gitee_user_endpoint"    required:"true"`
	GithubUserEndpoint string `json:"github_user_endpoint"   required:"true"`
}

func (cfg *omConfig) userEndpoint(userId, platform string) string {
	switch platform {
	case dp.Gitee:
		return cfg.GiteeUserEndpoint + userId

	case dp.Github:
		return cfg.GithubUserEndpoint + userId

	default:
		return ""
	}
}
