package dp

import "strings"

var config Config

func Init(cfg *Config, sv SigValidator) {
	config = *cfg
	sigValidator = sv
}

type Config struct {
	SupportedLanguages           []string `json:"supported_languages"`
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
	items := cfg.SupportedLanguages
	for i, s := range items {
		items[i] = strings.ToLower(s)
	}

	return nil
}

func (cfg *Config) isValidLanguage(v string) bool {
	v = strings.ToLower(v)

	for _, s := range cfg.SupportedLanguages {
		if v == s {
			return true
		}
	}

	return false
}
