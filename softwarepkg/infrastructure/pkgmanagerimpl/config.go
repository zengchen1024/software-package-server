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

type ExistingPkgsConfig struct {
	SigEndpoint  string                 `json:"sig_endpoint"      required:"true"`
	OrgOfPkgRepo string                 `json:"org_of_pkg_repo"   required:"true"`
	MetadataRepo MetadataRepoConfig     `json:"meta_data_repo"    required:"true"`
	DefaultInfo  ExistingPkgDefaultInfo `json:"default_info"      required:"true"`
}

func (cfg *ExistingPkgsConfig) setDefault() {
	if cfg.OrgOfPkgRepo == "" {
		cfg.OrgOfPkgRepo = "src-openeuler"
	}
}

type MetadataRepoConfig struct {
	Org    string `json:"org"      required:"true"`
	Repo   string `json:"repo"     required:"true"`
	Branch string `json:"branch"   required:"true"`
}

type ExistingPkgDefaultInfo struct {
	SpecURL        string `json:"spec_url"          required:"true"`
	SrcRPMURL      string `json:"src_rpm_url"       required:"true"`
	Platform       string `json:"platform"          required:"true"`
	RelevantPR     string `json:"relevant_pr"       required:"true"`
	ImporterName   string `json:"importer_name"     required:"true"`
	ImporterEmail  string `json:"importer_email"    required:"true"`
	ReasonToImport string `json:"reason_to_import"  required:"true"`
}

func (cfg *ExistingPkgDefaultInfo) toPkgBasicInfo() (info domain.SoftwarePkgBasicInfo, err error) {
	info.Phase = dp.PackagePhaseImported

	// importer
	importer := &info.Importer

	if importer.Account, err = dp.NewAccount(cfg.ImporterName); err != nil {
		return
	}

	if importer.Email, err = dp.NewEmail(cfg.ImporterEmail); err != nil {
		return
	}

	// pr
	if info.RelevantPR, err = dp.NewURL(cfg.RelevantPR); err != nil {
		return
	}

	// application
	app := &info.Application

	if app.PackagePlatform, err = dp.NewPackagePlatform(cfg.Platform); err != nil {
		return
	}

	if app.ReasonToImportPkg, err = dp.NewReasonToImportPkg(cfg.ReasonToImport); err != nil {
		return
	}

	// source code
	source := &app.SourceCode
	if source.SpecURL, err = dp.NewURL(cfg.SpecURL); err != nil {
		return
	}

	source.SrcRPMURL, err = dp.NewURL(cfg.SrcRPMURL)

	return
}
