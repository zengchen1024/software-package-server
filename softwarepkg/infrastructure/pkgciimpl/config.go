package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Config struct {
	CIRepo         CIRepo `json:"ci_repo"`
	WorkDir        string `json:"work_dir" required:"true"`
	InitScript     string `json:"init_script"`
	ClearScript    string `json:"clear_script"`
	DownloadScript string `json:"download_script"`
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

type CIRepo struct {
	// Org is the remote org of repo for CI
	Org string `json:"org"                      required:"true"`

	// Repo is the repo for CI. Suppose that the remote and local repo is the same name.
	Repo string `json:"repo"                    required:"true"`

	// Owner is the owner of local repo
	Owner string `json:"owner"                  required:"true"`
	Email string `json:"email"                  required:"true"`
	Token string `json:"token"                  required:"true"`

	// Link is the local repo address.
	Link string `json:"link"                    required:"true"`

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

func (cfg *CIRepo) cloneURL() string {
	return fmt.Sprintf(
		"https://%s:%s@%s",
		cfg.Owner, cfg.Token,
		strings.TrimPrefix(cfg.Link, "https://"),
	)
}
