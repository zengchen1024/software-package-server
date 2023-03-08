package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/utils"

	commonstl "github.com/opensourceways/software-package-server/common/controller"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

var (
	userInfoURL string
	client      utils.HttpClient
)

type Config struct {
	UserInfoURL string `json:"user_info_url" required:"true"`
}

const (
	keyCookie             = "_Y_G_"
	keyUserInfo           = "userinfo"
	headerPrivateToken    = "PRIVATE-TOKEN"
	errorBadRequestHeader = "bad_request_header"
	errorBadRequestCookie = "bad_request_cookie"
)

func Init(cfg *Config) {
	client = utils.NewHttpClient(3)
	userInfoURL = cfg.UserInfoURL
}

func CheckUser(ctx *gin.Context) {
	if code, err := checkUser(ctx); err != nil {
		commonstl.SendBadRequest(ctx, code, err)

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

func checkUser(ctx *gin.Context) (string, error) {
	t, err := token(ctx)
	if err != nil {
		return errorBadRequestHeader, err
	}

	c, err := cookie(ctx)
	if err != nil {
		return errorBadRequestCookie, err
	}

	v, err := getUserInfo(t, c)
	if err == nil {
		ctx.Set(keyUserInfo, v)
	}

	return "", err
}

func GetUser(ctx *gin.Context) (domain.User, error) {
	if v, exists := ctx.Get(keyUserInfo); exists {
		if u, ok := v.(domain.User); ok {
			return u, nil
		}
	}

	return domain.User{}, errors.New("no user info")
}

func token(ctx *gin.Context) (t string, err error) {
	if t = ctx.GetHeader(headerPrivateToken); len(t) == 0 {
		err = errors.New("invalid token")
	}

	return
}

func cookie(ctx *gin.Context) (c string, err error) {
	if c, err = ctx.Cookie(keyCookie); err != nil || len(c) == 0 {
		err = errors.New("invalid cookie")
	}

	return
}

func cookieHeader(cookie string) string {
	return keyCookie + "=" + cookie
}

type userInfoData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (d *userInfoData) toUser() (v domain.User, err error) {
	if v.Account, err = dp.NewAccount(d.Username); err != nil {
		return
	}

	v.Email, err = dp.NewEmail(d.Email)

	return
}

func getUserInfo(token, cookie string) (r domain.User, err error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("token", token)
	req.Header.Set("Cookie", cookieHeader(cookie))

	var result struct {
		Data userInfoData `json:"data"`
	}

	code, err := client.ForwardTo(req, &result)
	if err != nil {
		if code == http.StatusUnauthorized {
			err = errors.New("no login")
		}
	} else {
		r, err = result.Data.toUser()
	}

	return
}
