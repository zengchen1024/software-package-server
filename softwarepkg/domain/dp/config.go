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
	SupportedLanguages           []string `json:"supported_languages"       required:"true"`
	SupportedPlatforms           []string `json:"supported_platforms"       required:"true"`
	LocalPlatform                string   `json:"local_platform"            required:"true"`
	MaxLengthOfPackageName       int      `json:"max_length_of_pkg_name"`
	MaxLengthOfPackageDesc       int      `json:"max_length_of_pkg_desc"`
	MaxLengthOfReviewComment     int      `json:"max_length_of_review_comment"`
	MaxLengthOfReasonToImportPkg int      `json:"max_length_of_reason_to_import_pkg"`
}

func (cfg *Config) SetDefault() {
	if len(cfg.SupportedLanguages) == 0 {
		cfg.SupportedLanguages = []string{"chinese", "english"}
	}

	if cfg.MaxLengthOfPackageName <= 0 {
		cfg.MaxLengthOfPackageName = 50
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
	cfg.toLower(cfg.SupportedPlatforms)

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
	return cfg.has(strings.ToLower(v), cfg.SupportedPlatforms)
}

func (cfg *Config) has(v string, items []string) bool {
	for _, s := range items {
		if v == s {
			return true
		}
	}

	return false
}
