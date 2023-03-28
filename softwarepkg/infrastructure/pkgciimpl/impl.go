package pkgciimpl

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/pkgci"
)

var instance *pkgCIImpl

func Init(cfg *Config) error {
	data, err := ioutil.ReadFile(cfg.CommonTestBodyFile)
	if err != nil {
		return err
	}

	instance = &pkgCIImpl{
		cli:        utils.NewHttpClient(3),
		endpoint:   cfg.CIEndpoint,
		commonBody: bytes.NewReader(data),
	}

	return nil
}

func PkgCI() *pkgCIImpl {
	return instance
}

type softwarePkgInfo struct {
	PkgId     string `json:"pkg_id"`
	SpecURL   string `json:"spec_url"`
	SrcRPMURL string `json:"src_rpm_url"`
}

// pkgCIImpl
type pkgCIImpl struct {
	cli        utils.HttpClient
	endpoint   string
	commonBody io.Reader
}

func (impl *pkgCIImpl) SendTest(info *pkgci.SoftwarePkgInfo) error {
	v := softwarePkgInfo{
		PkgId:     info.PkgId,
		SpecURL:   info.SourceCode.SpecURL.URL(),
		SrcRPMURL: info.SourceCode.SrcRPMURL.URL(),
	}

	data, err := utils.JsonMarshal(v)
	if err != nil {
		return err
	}

	payload := testPayload{
		common:  impl.commonBody,
		pkgInfo: bytes.NewReader(data[1:]),
	}

	req, err := http.NewRequest(http.MethodPost, impl.endpoint, &payload)
	if err != nil {
		return err
	}

	_, _, err = impl.cli.Download(req)

	return err
}

type testPayload struct {
	common  io.Reader
	pkgInfo io.Reader
}

func (pkg *testPayload) Read(p []byte) (n int, err error) {
	if n, err = pkg.common.Read(p); n != len(p) {
		i := n
		n, err = pkg.pkgInfo.Read(p[n:])
		n += i
	}

	return
}
