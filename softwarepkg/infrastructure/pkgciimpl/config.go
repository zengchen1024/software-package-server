package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Config struct {
	CIRepo         CIRepo  `json:"ci_repo"`
	GitUser        GitUser `json:"git_user"`
	WorkDir        string  `json:"work_dir" required:"true"`
	InitScript     string  `json:"init_script"`
	ClearScript    string  `json:"clear_script"`
	DownloadScript string  `json:"download_script"`
}

func (cfg *Config) SetDefault() {
	cfg.CIRepo.setDefault()

	if cfg.InitScript == "" {
		cfg.InitScript = "/opt/app/init.sh"
	}

	if cfg.ClearScript == "" {
		cfg.ClearScript = "/opt/app/clear.sh"
	}

	if cfg.DownloadScript == "" {
		cfg.DownloadScript = "/opt/app/download.sh"
	}
}

type GitUser struct {
	User  string `json:"user"   required:"true"`
	Email string `json:"email"  required:"true"`
	Token string `json:"token"  required:"true"`
}

type CIRepo struct {
	Org         string `json:"org"              required:"true"`
	Repo        string `json:"repo"             required:"true"`
	Link        string `json:"link"             required:"true"`
	FileAddr    string `json:"file_addr"        required:"true"`
	MainBranch  string `json:"main_branch"`
	LFSFileAddr string `json:"lfs_file_addr"    required:"true"`
}

func (cfg *CIRepo) setDefault() {
	if cfg.MainBranch == "" {
		cfg.MainBranch = "master"
	}
}

func (cfg *CIRepo) fileAddr(pkgName dp.PackageName, fileName string, lfs bool) (dp.URL, error) {
	s := cfg.FileAddr
	if lfs {
		s = cfg.LFSFileAddr
	}

	return dp.NewURL(fmt.Sprintf(s, pkgName.PackageName(), fileName))
}

func (cfg *CIRepo) cloneURL(user *GitUser) string {
	return fmt.Sprintf(
		"https://%s:%s@%s",
		user.User, user.Token,
		strings.TrimPrefix(cfg.Link, "https://"),
	)
}
