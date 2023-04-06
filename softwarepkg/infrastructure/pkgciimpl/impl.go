package pkgciimpl

import (
	"bytes"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

var instance *pkgCIImpl

func Init(cfg *Config) {
	instance = &pkgCIImpl{
		cli:      utils.NewHttpClient(3),
		service:  cfg.CIService,
		endpoint: cfg.CIEndpoint,
	}
}

func PkgCI() *pkgCIImpl {
	return instance
}

type softwarePkgInfo struct {
	PkgId     string `json:"pkg_id"`
	PkgName   string `json:"pkg_name"`
	Service   string `json:"service"`
	SpecURL   string `json:"spec_url"`
	SrcRPMURL string `json:"src_rpm_url"`
}

// pkgCIImpl
type pkgCIImpl struct {
	cli      utils.HttpClient
	service  string
	endpoint string
}

func (impl *pkgCIImpl) SendTest(info *domain.SoftwarePkgBasicInfo) error {
	source := &info.Application.SourceCode
	v := softwarePkgInfo{
		PkgId:     info.Id,
		PkgName:   info.PkgName.PackageName(),
		Service:   impl.service,
		SpecURL:   source.SpecURL.URL(),
		SrcRPMURL: source.SrcRPMURL.URL(),
	}

	data, err := utils.JsonMarshal(v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, impl.endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}
