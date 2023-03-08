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
	headerPrivateToken    = "PRIVATE-TOKEN"
	cookieKey             = "_Y_G_"
	userKey               = "userinfo"
	errorBadRequestHeader = "bad_request_header"
	errorBadRequestCookie = "bad_request_cookie"
)

func Init(cfg *Config) {
	client = utils.NewHttpClient(3)
	userInfoURL = cfg.UserInfoURL
}

func cookieValue(cookie string) string {
	return cookieKey + "=" + cookie
}

func CheckUser(ctx *gin.Context) {
	t, err := token(ctx)
	if err != nil {
		commonstl.SendBadRequest(ctx, errorBadRequestHeader, err)
		ctx.Abort()

		return
	}

	c, err := cookie(ctx)
	if err != nil {
		commonstl.SendBadRequest(ctx, errorBadRequestCookie, err)
		ctx.Abort()

		return
	}

	v, err := getUserInfo(t, c)
	if err != nil {
		commonstl.SendBadRequest(ctx, "", err)
		ctx.Abort()

		return
	}

	ctx.Set(userKey, v)

	ctx.Next()
}

func GetUser(ctx *gin.Context) (*domain.User, error) {
	u, _ := ctx.Get(userKey)
	if userinfo, ok := u.(*domain.User); ok {
		return userinfo, nil
	}

	return nil, errors.New("no userinfo")
}

func token(ctx *gin.Context) (t string, err error) {
	if t = ctx.GetHeader(headerPrivateToken); len(t) == 0 {
		err = errors.New("invalid token")
	}

	return
}

func cookie(ctx *gin.Context) (c string, err error) {
	if c, err = ctx.Cookie(cookieKey); err != nil || len(c) == 0 {
		err = errors.New("invalid cookie")
	}

	return
}

func getUserInfo(t, c string) (userinfo *domain.User, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("token", t)
	req.Header.Set("Cookie", cookieValue(c))

	var result = struct {
		Data struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"data"`
	}{}

	code, _ := client.ForwardTo(req, &result)
	if code == http.StatusUnauthorized {
		err = errors.New("no login")

		return
	}

	userinfo = new(domain.User)
	userinfo.Account, err = dp.NewAccount(result.Data.Username)
	if err != nil {
		return
	}

	userinfo.Email, err = dp.NewEmail(result.Data.Email)
	if err != nil {
		return
	}

	return
}
