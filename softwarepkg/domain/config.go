package domain

import (
	"errors"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	pkgModificationSig       = "sig"
	pkgModificationCode      = "code"
	pkgModificationPkgName   = "pkg_name"
	pkgModificationPkgDesc   = "pkg_desc"
	pkgModificationUpstream  = "upstream"
	pkgModificationPkgReason = "pkg_reason"
	pkgModificationCommitter = "committer"
)

var (
	timeoutOfCI      int64
	commonCheckItems []CheckItem
)

func Init(cfg *Config, m maintainer, ci pkgCI) {
	timeoutOfCI = cfg.CITimeout
	commonCheckItems = cfg.checkItems

	ciInstance = ci
	maintainerInstance = m
}

func pkgModification(v string) (string, error) {
	v = strings.ToLower(v)

	switch v {
	case pkgModificationSig:
	case pkgModificationCode:
	case pkgModificationPkgName:
	case pkgModificationPkgDesc:
	case pkgModificationUpstream:
	case pkgModificationPkgReason:
	case pkgModificationCommitter:
	default:
		return "", errors.New("invalid pkg modification")
	}

	return v, nil
}

type Config struct {
	CITimeout  int64             `json:"ci_timeout"`
	CheckItems []checkItemConfig `json:"check_items"`

	checkItems []CheckItem
}

func (cfg *Config) SetDefault() {
	if cfg.CITimeout <= 0 {
		cfg.CITimeout = 3 * 3600 // 3 hours
	}
}

func (cfg *Config) Validate() (err error) {
	v := make([]CheckItem, len(cfg.CheckItems))

	for i := range cfg.CheckItems {
		if v[i], err = cfg.CheckItems[i].toCheckItem(); err != nil {
			return
		}
	}

	cfg.checkItems = v

	return
}

type checkItemConfig struct {
	Id    string `json:"id"     required:"true"`
	Name  string `json:"name"   required:"true"`
	Desc  string `json:"desc"   required:"true"`
	Owner string `json:"owner"  required:"true"`

	// This check item should be checked again when the relevant modifications happened.
	Modifications []string `json:"modifications" required:"true"`

	// If true, keep the review record of reviewer, otherwise clear all the records about
	// this item when relevant modifications happened.
	// For example, the review about the item whether the user aggreed to
	// to be committer of the pkg should be kept.
	Keep bool `json:"keep"`

	// If true, only the owner can review this item else anyone can review.
	// For example, onlye sig maintainer can determine whether the sig of pkg is correct.
	OnlyOwner bool `json:""`
}

func (cfg *checkItemConfig) toCheckItem() (item CheckItem, err error) {
	if item.Owner, err = dp.NewCommunityRole(cfg.Owner); err != nil {
		return
	}

	ms := make([]string, len(cfg.Modifications))
	for i, v := range cfg.Modifications {
		if ms[i], err = pkgModification(v); err != nil {
			return
		}
	}
	item.Modifications = ms

	item.Id = cfg.Id
	item.Name = cfg.Name
	item.Desc = cfg.Desc
	item.Keep = cfg.Keep
	item.OnlyOwner = cfg.OnlyOwner

	return
}
