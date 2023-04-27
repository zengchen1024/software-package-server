package pkgciimpl

import (
	"fmt"
	"strings"
)

type Config struct {
	WorkDir      string  `json:"work_dir"       required:"true"`
	GitUser      GitUser `json:"user"           required:"true"`
	CIRepo       CIRepo  `json:"ci_repo"        required:"true"`
	PRScript     string  `json:"pr_script"      required:"true"`
	CloneScript  string  `json:"clone_script"   required:"true"`
	CIComment    string  `json:"ci_comment"     required:"true"`
	CIService    string  `json:"ci_service"     required:"true"`
	TargetBranch string  `json:"target_branch"  required:"true"`
}

type GitUser struct {
	User  string `json:"user"   required:"true"`
	Email string `json:"email"  required:"true"`
	Token string `json:"token"  required:"true"`
}

type CIRepo struct {
	Org  string `json:"org"     required:"true"`
	Repo string `json:"repo"    required:"true"`
	Link string `json:"link"    required:"true"`
}

func (cfg *CIRepo) cloneURL(user *GitUser) string {
	return fmt.Sprintf(
		"https://%s:%s@%s",
		user.User, user.Token,
		strings.TrimPrefix(cfg.Link, "https://"),
	)
}

func (cfg *Config) SetDefault() {
	if cfg.PRScript == "" {
		cfg.PRScript = "/opt/app/pull_request.sh"
	}

	if cfg.CloneScript == "" {
		cfg.CloneScript = "/opt/app/clone_repo.sh"
	}

	if cfg.TargetBranch == "" {
		cfg.TargetBranch = "master"
	}

	if cfg.CIComment == "" {
		cfg.CIComment = "/retest"
	}
}
