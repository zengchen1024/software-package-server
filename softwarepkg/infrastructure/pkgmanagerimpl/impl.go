package pkgmanagerimpl

import (
	"errors"
	"fmt"
	"net/http"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	libutils "github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

var instance *service

func Init(cfg *Config) error {
	info := &cfg.ExistingPkgs

	v, err := info.DefaultInfo.toPkgBasicInfo()
	if err != nil {
		return err
	}

	instance = &service{
		cli:              client.NewClient(cfg.Token()),
		httpCli:          libutils.NewHttpClient(3),
		defaultPkg:       v,
		orgOfPkgRepo:     info.OrgOfPkgRepo,
		metaDataEndpoint: info.MetaDataEndpoint,
	}

	return nil
}

func Instance() *service {
	return instance
}

type pkgMetaData struct {
	Description string `json:"description"`
	SigName     string `json:"sig_name"`
}

type service struct {
	cli              client.Client
	httpCli          libutils.HttpClient
	defaultPkg       domain.SoftwarePkgBasicInfo
	orgOfPkgRepo     string
	metaDataEndpoint string
}

func (s *service) IsPkgExisted(pkg dp.PackageName) bool {
	_, err := s.cli.GetRepo(s.orgOfPkgRepo, pkg.PackageName())

	return err == nil
}

func (s *service) GetPkg(name dp.PackageName) (info domain.SoftwarePkgBasicInfo, err error) {
	repo, err := s.cli.GetRepo(s.orgOfPkgRepo, name.PackageName())
	if err != nil {
		return
	}

	meta, err := s.getPkgMetaData(name)
	if err != nil {
		return
	}

	return s.toPkgBasicInfo(name, &repo, &meta)
}

func (s *service) toPkgBasicInfo(
	name dp.PackageName, repo *sdk.Project, meta *pkgMetaData,
) (info domain.SoftwarePkgBasicInfo, err error) {

	info = s.defaultPkg

	info.PkgName = name
	info.AppliedAt = utils.Now()

	url, err := dp.NewURL(repo.GetHtmlUrl())
	if err != nil {
		return
	}

	info.RepoLink = url
	info.RelevantPR = url

	app := &info.Application
	app.SourceCode.SrcRPMURL = url
	app.SourceCode.SpecURL = url

	desc := repo.Description
	if desc == "" {
		desc = fmt.Sprintf("importing software package: %s", name.PackageName())
	}

	if app.PackageDesc, err = dp.NewPackageDesc(desc); err != nil {
		return
	}

	if app.ImportingPkgSig, err = dp.NewImportingPkgSig(meta.SigName); err != nil {
		return
	}

	return
}

func (s *service) getPkgMetaData(name dp.PackageName) (r pkgMetaData, err error) {
	req, err := http.NewRequest(
		http.MethodGet, s.metaDataEndpoint+name.PackageName(), nil,
	)
	if err != nil {
		return
	}

	var v struct {
		Data []pkgMetaData `json:"data"`
	}

	if _, err = s.httpCli.ForwardTo(req, &v); err != nil {
		return
	}

	if len(v.Data) == 0 {
		err = errors.New("pkg meta data is not found")
	} else {
		r = v.Data[0]
	}

	return
}
