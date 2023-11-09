package pkgmanagerimpl

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type Config struct {
	ExistingPkgs ExistingPkgsConfig `json:"existing_pkgs"  required:"true"`
	AccessToken  string             `json:"access_token"   required:"true"`
}

func (cfg *Config) SetDefault() {
	cfg.ExistingPkgs.setDefault()
}

func (cfg *Config) Token() func() []byte {
	return func() []byte {
		return []byte(cfg.AccessToken)
	}
}

// ExistingPkgsConfig
type ExistingPkgsConfig struct {
	DefaultInfo      ExistingPkgDefaultInfo `json:"default_info"        required:"true"`
	MetaDataRepo     metaDataRepo           `json:"meta_data_repo"      required:"true"`
	OrgOfPkgRepo     string                 `json:"org_of_pkg_repo"     required:"true"`
	MetaDataEndpoint string                 `json:"meta_data_endpoint"  required:"true"`
}

func (cfg *ExistingPkgsConfig) setDefault() {
	if cfg.OrgOfPkgRepo == "" {
		cfg.OrgOfPkgRepo = "src-openeuler"
	}

	cfg.MetaDataRepo.setDefault()
}

// metaDataRepo
type metaDataRepo struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

func (cfg *metaDataRepo) setDefault() {
	if cfg.Owner == "" {
		cfg.Owner = "openeuler"
	}

	if cfg.Repo == "" {
		cfg.Repo = "community"
	}

	if cfg.Branch == "" {
		cfg.Branch = "master"
	}
}

// ExistingPkgDefaultInfo
type ExistingPkgDefaultInfo struct {
	Platform       string `json:"platform"          required:"true"`
	ImporterName   string `json:"importer_name"     required:"true"`
	ReasonToImport string `json:"reason_to_import"  required:"true"`
}

func (cfg *ExistingPkgDefaultInfo) toPkgBasicInfo() (info domain.SoftwarePkg, err error) {
	info.Phase = dp.PackagePhaseImported

	// importer
	if info.Importer, err = dp.NewAccount(cfg.ImporterName); err != nil {
		return
	}

	// application
	if info.Repo.Platform, err = dp.NewPackagePlatform(cfg.Platform); err != nil {
		return
	}

	if info.Basic.Reason, err = dp.NewReasonToImportPkg(cfg.ReasonToImport); err != nil {
		return
	}

	return
}
