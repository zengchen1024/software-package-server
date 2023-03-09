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
	isMaintainer := func(v []string) bool {
		for _, t := range v {
			if t == maintainers {
				return true
			}
		}

		return false
	}

	sig = strings.ToLower(sig)

	data := s.Data
	for i := range data {
		if sig == strings.ToLower(data[i].Sig) {
			return isMaintainer(data[i].Type)
		}
	}

	return false
}

func NewMaintainerImpl(cfg *Config) *maintainerImpl {
	return &maintainerImpl{
		cli:           utils.NewHttpClient(3),
		permissionURL: cfg.PermissionURL + "&user=",
	}
}

type maintainerImpl struct {
	cli           utils.HttpClient
	permissionURL string
}

func (impl *maintainerImpl) baseUrl(user string) string {
	return impl.permissionURL + user
}

func (impl *maintainerImpl) HasPermission(info *domain.SoftwarePkgBasicInfo, user dp.Account) (
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

	return res.hasPermission(
		info.Application.ImportingPkgSig.ImportingPkgSig(),
	), nil
}
