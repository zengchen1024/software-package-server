package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	EcoPkgSig                     string `json:"ecopkg_sig"`
	MinNumApprovedByTC            int    `json:"min_num_approved_by_tc"`
	MinNumApprovedBySigMaintainer int    `json:"min_num_approved_by_sig_maintainer"`
}

func (cfg *Config) SetDefault() {
	if cfg.EcoPkgSig == "" {
		cfg.EcoPkgSig = "ecopkg"
	}

	if cfg.MinNumApprovedByTC <= 0 {
		cfg.MinNumApprovedByTC = 1
	}

	if cfg.MinNumApprovedBySigMaintainer <= 0 {
		cfg.MinNumApprovedBySigMaintainer = 2
	}
}
