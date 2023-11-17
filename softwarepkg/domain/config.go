package domain

var config Config

func Init(cfg *Config, m maintainer, ci pkgCI) {
	config = *cfg
	ciInstance = ci
	maintainerInstance = m
}

type Config struct {
	EcopkgSig                     string `json:"ecopkg_sig"`
	CITimeout                     int64  `json:"ci_timeout"`
	MinNumApprovedByTC            int    `json:"min_num_approved_by_tc"`
	MinNumApprovedBySigMaintainer int    `json:"min_num_approved_by_sig_maintainer"`
}

func (cfg *Config) SetDefault() {
	if cfg.EcopkgSig == "" {
		cfg.EcopkgSig = "ecopkg"
	}

	if cfg.CITimeout <= 0 {
		cfg.CITimeout = 3 * 3600 // 3 hours
	}

	if cfg.MinNumApprovedByTC <= 0 {
		cfg.MinNumApprovedByTC = 1
	}

	if cfg.MinNumApprovedBySigMaintainer <= 0 {
		cfg.MinNumApprovedBySigMaintainer = 2
	}
}
