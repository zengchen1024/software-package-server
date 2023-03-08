package maintainerimpl

import (
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	maintainers = "maintainers"
)

type sigPermission struct {
	Data []struct {
		Sig  string   `json:"sig"`
		Type []string `json:"type"`
	} `json:"data"`
}

func (s sigPermission) hasPermission(sig string) bool {
	for _, v := range s.Data {
		if strings.EqualFold(sig, v.Sig) {
			for _, t := range v.Type {
				if t == maintainers {
					return true
				}
			}
		}
	}

	return false
}

func NewMaintainerImpl(cfg *Config) maintainerImpl {
	cfg.PermissionURL = cfg.PermissionURL + "&user="
	return maintainerImpl{
		cfg: *cfg,
		cli: utils.NewHttpClient(3),
	}
}

type maintainerImpl struct {
	cfg Config
	cli utils.HttpClient
}

func (impl maintainerImpl) baseUrl(user string) string {
	return impl.cfg.PermissionURL + user
}

func (impl maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user dp.Account) (
	bool, error,
) {
	req, err := http.NewRequest(http.MethodGet, impl.baseUrl(user.Account()), nil)
	if err != nil {
		return false, err
	}

	var res sigPermission
	if _, err = impl.cli.ForwardTo(req, &res); err != nil {
		return false, err
	}

	return res.hasPermission(info.Application.ImportingPkgSig.ImportingPkgSig()), nil
}
