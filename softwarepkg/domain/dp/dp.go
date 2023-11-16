package dp

import (
	"errors"
	"net/url"
	"path"
	"regexp"

	libutil "github.com/opensourceways/server-common-lib/utils"
)

var (
	reName = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || !reName.MatchString(v) {
		return nil, errors.New("invalid account")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

func IsSameAccount(a, b Account) bool {
	return a != nil && b != nil && a.Account() == b.Account()
}

// URL
type URL interface {
	URL() string
	FileName() string
}

func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.ParseRequestURI(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpURL(v), nil
}

type dpURL string

func (v dpURL) URL() string {
	return string(v)
}

func (v dpURL) FileName() string {
	r, err := url.ParseRequestURI(string(v))
	if err != nil {
		return ""
	}

	return path.Base(r.Path)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v == "" || !libutil.IsValidEmail(v) {
		return nil, errors.New("invalid email")
	}

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}
