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
	v, err := cfg.ExistingPkgs.DefaultInfo.toPkgBasicInfo()
	if err != nil {
		return err
	}

	instance = &service{
		cli:        client.NewClient(cfg.Token()),
		cfg:        cfg.ExistingPkgs,
		libcli:     libutils.NewHttpClient(3),
		defaultPkg: v,
	}

	return nil
}

func Instance() *service {
	return instance
}

type pkgMetaData struct {
	Data []metaData `json:"data"`
}

type metaData struct {
	MailingList string `json:"mailing_list"`
	Description string `json:"description"`
	SigName     string `json:"sig_name"`
}

type service struct {
	cli        client.Client
	cfg        ExistingPkgsConfig
	libcli     libutils.HttpClient
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

	meta, err := s.getPkgMetaData(name)
	if err != nil {
		return
	}

	return s.toPkgBasicInfo(name, &repo, &meta.Data[0])
}

func (s *service) toPkgBasicInfo(
	name dp.PackageName, repo *sdk.Project, meta *metaData,
) (info domain.SoftwarePkgBasicInfo, err error) {

	info = s.defaultPkg

	info.PkgName = name
	info.AppliedAt = utils.Now()

	url, err := dp.NewURL(repo.GetUrl())
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

	req, err := http.NewRequest(http.MethodGet, s.cfg.MetaDataEndpoint+name.PackageName(), nil)
	if err != nil {
		return pkgMetaData{}, err
	}

	if _, err = s.libcli.ForwardTo(req, &r); err != nil {
		return pkgMetaData{}, err
	}

	if len(r.Data) == 0 {
		err = errors.New("not find pkg meta data")
	}

	return
}
