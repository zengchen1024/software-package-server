package useradapterimpl

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type omTokenReq struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

type omTokenResp struct {
	Msg    string `json:"msg"`
	Token  string `json:"token"`
	Status int    `json:"status"`
}

type omUserInfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Email      string           `json:"email"`
		Username   string           `json:"username"`
		Identities []identitiesResp `json:"identities"`
	} `json:"data"`
}

func (resp *omUserInfoResp) toUser() (user domain.User, err error) {
	for _, item := range resp.Data.Identities {
		switch item.Identity {
		case dp.Gitee:
			user.GiteeID = item.LoginName

		case dp.Github:
			user.GithubID = item.LoginName

		default:
		}
	}

	if user.Account, err = dp.NewAccount(resp.Data.Username); err != nil {
		return
	}

	user.Email, err = dp.NewEmail(resp.Data.Email)

	return
}

type identitiesResp struct {
	LoginName string `json:"login_name"`
	Identity  string `json:"identity"`
}

type omClient struct {
	config omConfig
}

func (om *omClient) getToken() (string, error) {
	payload, err := utils.JsonMarshal(omTokenReq{
		GrantType: "token",
		AppId:     om.config.AppId,
		AppSecret: om.config.AppSecret,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost, om.config.TokenEndpoint, bytes.NewBuffer(payload),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	v := new(omTokenResp)
	cli := utils.NewHttpClient(3)

	if _, err = cli.ForwardTo(req, v); err != nil {
		return "", err
	}

	if v.Status != 200 {
		return "", errors.New(v.Msg)
	}

	return v.Token, nil
}

func (om *omClient) getUserInfo(userId, platform string) (user domain.User, err error) {
	token, err := om.getToken()
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, om.config.userEndpoint(userId, platform), nil)
	if err != nil {
		return
	}

	req.Header.Set("token", token)

	v := new(omUserInfoResp)
	cli := utils.NewHttpClient(3)

	if _, err = cli.ForwardTo(req, v); err != nil {
		return
	}

	if v.Code != 200 {
		err = errors.New(v.Msg)

		return
	}

	return v.toUser()
}
