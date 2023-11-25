package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Config struct {
	WorkDir      string  `json:"work_dir"       required:"true"`
	GitUser      GitUser `json:"user"           required:"true"`
	PRScript     string  `json:"pr_script"      required:"true"`
	CloneScript  string  `json:"clone_script"   required:"true"`
	CIRepo       CIRepo  `json:"ci_repo"        required:"true"`
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
	Org         string `json:"org"     required:"true"`
	Repo        string `json:"repo"    required:"true"`
	Link        string `json:"link"    required:"true"`
	FileAddr    string `json:"file_addr"    required:"true"`
	LFSFileAddr string `json:"lfs_file_addr"    required:"true"`
}

func (cfg *CIRepo) fileAddr(name string, lfs bool) (dp.URL, error) {
	if lfs {
		return dp.NewURL(fmt.Sprintf(cfg.LFSFileAddr, name))
	}

	return dp.NewURL(fmt.Sprintf(cfg.LFSFileAddr, name))
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
