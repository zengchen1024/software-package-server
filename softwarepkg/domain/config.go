package domain

var config Config

func Init(cfg *Config, m maintainer) {
	config = *cfg
	maintainerInstance = m
}

type Config struct {
	EcopkgSig                     string `json:"ecopkg_sig"`
	MinNumApprovedByTC            int    `json:"min_num_approved_by_tc"`
	MinNumApprovedBySigMaintainer int    `json:"min_num_approved_by_sig_maintainer"`
}

func (cfg *Config) SetDefault() {
	if cfg.EcopkgSig == "" {
		cfg.EcopkgSig = "ecopkg"
	}

	if cfg.MinNumApprovedByTC <= 0 {
		cfg.MinNumApprovedByTC = 1
	}

	if cfg.MinNumApprovedBySigMaintainer <= 0 {
		cfg.MinNumApprovedBySigMaintainer = 2
	}
}
