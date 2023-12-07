package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/allerror"
	commonstl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	keyCookie                  = "_Y_G_"
	keyUserInfo                = "user_info"
	headerPrivateToken         = "PRIVATE-TOKEN"
	errorBadRequestHaventLogin = "bad_request_havent_login"
	errorNoLoginMsg            = "no login"
)

var instance *userCheckingMiddleware

type Config struct {
	UserInfoURL string `json:"user_info_url" required:"true"`
}

func Init(cfg *Config) {
	instance = &userCheckingMiddleware{
		client:      utils.NewHttpClient(3),
		userInfoURL: cfg.UserInfoURL,
	}
}

func UserChecking() *userCheckingMiddleware {
	return instance
}

type userCheckingMiddleware struct {
	userInfoURL string
	client      utils.HttpClient
}

func (m *userCheckingMiddleware) FetchUser(ctx *gin.Context) (domain.User, error) {
	if v, exists := ctx.Get(keyUserInfo); exists {
		if u, ok := v.(domain.User); ok {
			return u, nil
		}
	}

	return domain.User{}, errors.New("no user info")
}

func (m *userCheckingMiddleware) CheckUser(ctx *gin.Context) {
	if err := m.doCheck(ctx); err != nil {
		commonstl.SendError(ctx, err)

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

func (m *userCheckingMiddleware) doCheck(ctx *gin.Context) error {
	t := m.token(ctx)
	if t == "" {
		return allerror.New(errorBadRequestHaventLogin, errorNoLoginMsg)
	}

	c := m.cookie(ctx)
	if c == "" {
		return allerror.New(errorBadRequestHaventLogin, errorNoLoginMsg)
	}

	v, err := m.getUserInfo(t, c)
	if err == nil {
		ctx.Set(keyUserInfo, v)
	}

	return err
}

func (m *userCheckingMiddleware) token(ctx *gin.Context) string {
	return ctx.GetHeader(headerPrivateToken)
}

func (m *userCheckingMiddleware) cookie(ctx *gin.Context) string {
	c, _ := ctx.Cookie(keyCookie)

	return c
}

func (m *userCheckingMiddleware) cookieHeader(cookie string) string {
	return keyCookie + "=" + cookie
}

func (m *userCheckingMiddleware) getUserInfo(token, cookie string) (
	r domain.User, err error,
) {
	req, err := http.NewRequest(http.MethodGet, m.userInfoURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("token", token)
	req.Header.Set("Cookie", m.cookieHeader(cookie))

	var result struct {
		Data userInfoData `json:"data"`
	}

	status, err := m.client.ForwardTo(req, &result)
	if err != nil {
		if status == http.StatusUnauthorized {
			err = allerror.New(errorBadRequestHaventLogin, errorNoLoginMsg)
		}
	} else {
		r, err = result.Data.toUser()
	}

	return
}

type userInfoData struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Identities []struct {
		LoginName string `json:"login_name"`
		Identity  string `json:"identity"`
	} `json:"identities"`
}

func (d *userInfoData) toUser() (v domain.User, err error) {
	if v.Account, err = dp.NewAccount(d.Username); err != nil {
		return
	}

	v.Email, err = dp.NewEmail(d.Email)

	for _, identity := range d.Identities {
		switch identity.Identity {
		case dp.Gitee:
			v.GiteeID = identity.LoginName

		case dp.Github:
			v.GithubID = identity.LoginName

		default:
		}
	}

	return
}
