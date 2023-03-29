package pkgmanagerimpl

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

var instance *service

func Init(cfg *Config) error {
	v, err := cfg.ExistingPkgs.DefaultInfo.toPkgBasicInfo()
	if err != nil {
		return err
	}

	instance = &service{
		cli:        client.NewClient(cfg.Token()),
		cfg:        cfg.ExistingPkgs,
		defaultPkg: v,
	}

	return nil
}

func Instance() *service {
	return instance
}

type pkgMetaData struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type service struct {
	cli        client.Client
	cfg        ExistingPkgsConfig
	defaultPkg domain.SoftwarePkgBasicInfo
}

func (s *service) IsPkgExisted(pkg dp.PackageName) bool {
	_, err := s.cli.GetRepo(s.cfg.OrgOfPkgRepo, pkg.PackageName())

	return err == nil
}

func (s *service) GetPkg(name dp.PackageName) (info domain.SoftwarePkgBasicInfo, err error) {
	repo, err := s.cli.GetRepo(s.cfg.OrgOfPkgRepo, name.PackageName())
	if err != nil {
		return
	}

	sig, err := s.fetchSig(name)
	if err != nil {
		return
	}

	meta, err := s.getPkgMetaData(name, sig)
	if err != nil {
		return
	}

	return s.toPkgBasicInfo(name, &repo, &meta, sig)
}

func (s *service) toPkgBasicInfo(
	name dp.PackageName, repo *sdk.Project, meta *pkgMetaData, sig string,
) (info domain.SoftwarePkgBasicInfo, err error) {
	info = s.defaultPkg

	info.PkgName = name
	info.AppliedAt = utils.Now()

	if info.RepoLink, err = dp.NewURL(repo.GetUrl()); err != nil {
		return
	}

	app := &info.Application

	desc := repo.Description
	if desc == "" {
		desc = fmt.Sprintf("importing software package:%s", name.PackageName())
	}
	if app.PackageDesc, err = dp.NewPackageDesc(desc); err != nil {
		return
	}

	if app.ImportingPkgSig, err = dp.NewImportingPkgSig(sig); err != nil {
		return
	}

	return
}

func (s *service) fetchSig(pkg dp.PackageName) (sig string, err error) {
	// TODO
	return "", nil
}

func (s *service) getPkgMetaData(name dp.PackageName, sig string) (r pkgMetaData, err error) {
	meta := &s.cfg.MetadataRepo

	str := name.PackageName()
	path := fmt.Sprintf("sig/%s/src-openeuler/%s/%s.yaml", sig, string(str[0]), str)

	v, err := s.cli.GetPathContent(meta.Org, meta.Repo, path, meta.Branch)
	if err == nil {
		err = decodeYamlFile(v.Content, &r)
	}

	return
}

func decodeYamlFile(content string, v interface{}) error {
	c, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(c, v)
}
