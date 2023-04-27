package pkgciimpl

import (
	"fmt"
	"strings"
)

type Config struct {
	WorkDir string `json:"work_dir"           required:"true"`

	GitUser
	CIRepo

	CIScript    string `json:"ci_script"          required:"true"`
	CloneScript string `json:"clone_script"       required:"true"`

	Comment         string `json:"comment"            required:"true"`
	CIService       string `json:"ci_service"         required:"true"`
	CreateBranch    string `json:"create_branch"      required:"true"`
	CreateCIPRToken string `json:"create_ci_pr_token" required:"true"`
}

type GitUser struct {
	User  string `json:"user"   required:"true"`
	Email string `json:"email"  required:"true"`
	Token string `json:"token"  required:"true"`
}

type CIRepo struct {
	Org  string `json:"ci_org"   required:"true"`
	Repo string `json:"ci_repo"  required:"true"`
	Link string `json:"link"     required:"true"`
}

func (cfg *CIRepo) cloneURL(user *GitUser) string {
	return fmt.Sprintf(
		"https://%s:%s@%s",
		user.User, user.Token,
		strings.TrimPrefix(cfg.Link, "https://"),
	)
}

func (cfg *Config) SetDefault() {
	if cfg.CIScript == "" {
		cfg.CIScript = "/opt/app/pull_request.sh"
	}

	if cfg.CreateBranch == "" {
		cfg.CreateBranch = "master"
	}

	if cfg.Comment == "" {
		cfg.Comment = "/retest"
	}
}
