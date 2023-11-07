package dp

import (
	"errors"
	"strings"
)

var config Config

func Init(cfg *Config, sv SigValidator) {
	config = *cfg
	sigValidator = sv
}

type Config struct {
	// map platform -> community address. such as gitee -> https://gitee.com/src-openeuler/
	SupportedPlatforms           map[string]string `json:"supported_platforms"       required:"true"`
	SupportedLanguages           []string          `json:"supported_languages"       required:"true"`
	LocalPlatform                string            `json:"local_platform"            required:"true"`
	MaxLengthOfPackageName       int               `json:"max_length_of_pkg_name"`
	MaxLengthOfPackageDesc       int               `json:"max_length_of_pkg_desc"`
	MaxLengthOfReviewComment     int               `json:"max_length_of_review_comment"`
	MaxLengthOfReasonToImportPkg int               `json:"max_length_of_reason_to_import_pkg"`
}

func (cfg *Config) SetDefault() {
	if len(cfg.SupportedLanguages) == 0 {
		cfg.SupportedLanguages = []string{"chinese", "english"}
	}

	if cfg.MaxLengthOfPackageName <= 0 {
		cfg.MaxLengthOfPackageName = 60
	}

	if cfg.MaxLengthOfPackageDesc <= 0 {
		cfg.MaxLengthOfPackageDesc = 1000
	}

	if cfg.MaxLengthOfReviewComment <= 0 {
		cfg.MaxLengthOfReviewComment = 500
	}

	if cfg.MaxLengthOfReasonToImportPkg <= 0 {
		cfg.MaxLengthOfReasonToImportPkg = 1000
	}
}

func (cfg *Config) Validate() error {
	cfg.toLower(cfg.SupportedLanguages)

	m := map[string]string{}
	for k, v := range cfg.SupportedPlatforms {
		if !strings.HasSuffix(v, "/") {
			v += "/"
		}

		m[strings.ToLower(k)] = v
	}

	cfg.SupportedPlatforms = m

	if !cfg.isValidPlatform(cfg.LocalPlatform) {
		return errors.New("unkonw local platform")
	}

	return nil
}

func (cfg *Config) toLower(items []string) {
	for i, s := range items {
		items[i] = strings.ToLower(s)
	}
}

func (cfg *Config) isValidLanguage(v string) bool {
	return cfg.has(strings.ToLower(v), cfg.SupportedLanguages)
}

func (cfg *Config) isValidPlatform(v string) bool {
	return cfg.SupportedPlatforms[strings.ToLower(v)] != ""
}

func (cfg *Config) has(v string, items []string) bool {
	for _, s := range items {
		if v == s {
			return true
		}
	}

	return false
}
