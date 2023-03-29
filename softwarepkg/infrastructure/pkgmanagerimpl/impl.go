package pkgmanagerimpl

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

var (
	instance *service
	base     *BaseConfig
)

func Init(cfg *Config) {
	instance = &service{
		cli: client.NewClient(cfg.Token()),
		org: cfg.Org,
	}

	base = &cfg.Base
}

func Instance() *service {
	return instance
}

type repository struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type service struct {
	cli client.Client
	org string
}

func (s *service) IsPkgExisted(pkg dp.PackageName) bool {
	_, err := s.cli.GetRepo(s.org, pkg.PackageName())

	return err == nil
}

func (s *service) GetPkg(pkg dp.PackageName) (domain.SoftwarePkgBasicInfo, error) {
	v, err := s.cli.GetRepo(s.org, pkg.PackageName())
	if err != nil {
		return domain.SoftwarePkgBasicInfo{}, err
	}

	path, sig, err := s.repoDetailPath(pkg)
	if err != nil {
		return domain.SoftwarePkgBasicInfo{}, err
	}

	var repo repository
	if err = s.pathContent(path, &repo); err != nil {
		return domain.SoftwarePkgBasicInfo{}, err
	}

	return s.toPkgBasicInfo(pkg, v, repo, sig)
}

func (s *service) toPkgBasicInfo(
	pkg dp.PackageName, v sdk.Project, repo repository, sig string,
) (info domain.SoftwarePkgBasicInfo, err error) {
	info.PkgName = pkg
	importer := &info.Importer

	if importer.Account, err = dp.NewAccount("software-pkg-robot"); err != nil {
		return
	}

	if importer.Email, err = dp.NewEmail("software@openeuler.org"); err != nil {
		return
	}

	if info.RepoLink, err = dp.NewURL(v.GetUrl()); err != nil {
		return
	}

	info.Phase = dp.PackagePhaseImported
	info.AppliedAt = utils.Now()

	app := &info.Application
	if app.PackagePlatform, err = dp.NewPackagePlatform("gitee"); err != nil {
		return
	}

	desc := pkg.PackageName()
	if len(repo.Description) > 0 {
		desc = repo.Description
	}

	if app.PackageDesc, err = dp.NewPackageDesc(desc); err != nil {
		return
	}

	if app.ImportingPkgSig, err = dp.NewImportingPkgSig(sig); err != nil {
		return
	}

	if app.ReasonToImportPkg, err = dp.NewReasonToImportPkg(pkg.PackageName()); err != nil {
		return
	}

	source := &app.SourceCode
	if source.SpecURL, err = dp.NewURL(v.GetUrl()); err != nil {
		return
	}

	source.SrcRPMURL, err = dp.NewURL(v.GetUrl())

	return
}

func (s *service) repoDetailPath(pkg dp.PackageName) (
	path string, sig string, err error,
) {
	v, err := s.cli.GetDirectoryTree(base.Org, base.Repo, base.Branch, 1)
	if err != nil {
		return
	}

	for _, t := range v.Tree {
		patharr := strings.Split(t.Path, "/")
		if len(patharr) != 5 {
			continue
		}
		file := fmt.Sprintf("%s.yaml", pkg.PackageName())
		if patharr[0] == "sig" && patharr[2] == s.org && patharr[4] == file {
			path = t.Path
			sig = patharr[1]

			return
		}
	}

	return "", "", errors.New("no path")
}

func (s *service) pathContent(path string, res interface{}) error {
	v, err := s.cli.GetPathContent(base.Org, base.Repo, path, base.Branch)
	if err != nil {
		return err
	}

	return decodeYamlFile(v.Content, res)
}

func decodeYamlFile(content string, v interface{}) error {
	c, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(c, v)
}
