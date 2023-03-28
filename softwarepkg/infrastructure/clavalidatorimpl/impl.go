package clavalidatorimpl

import (
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var instance *service

func Init(cfg *Config) {
	instance = &service{
		url: cfg.CheckURL + "?email=",
		cli: utils.NewHttpClient(3),
	}
}

func Instance() *service {
	return instance
}

type service struct {
	url string
	cli utils.HttpClient
}

type signedResult struct {
	Data struct {
		Signed bool `json:"signed"`
	} `json:"data"`
}

func (s *service) HasSignedCLA(email dp.Email) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, s.url+email.Email(), nil)
	if err != nil {
		return false, err
	}

	v := new(signedResult)

	if _, err = s.cli.ForwardTo(req, v); err != nil {
		return false, err
	}

	return v.Data.Signed, nil
}
