package maintainerimpl

import "time"

type Config struct {
	ConfigForPermission

	ReadURL string `json:"read_url" required:"true"`

	// Interval the unit is hour
	Interval int `json:"interval" required:"true"`
}

func (cfg *Config) SetDefault() {
	cfg.setDefault()
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Duration(cfg.Interval) * time.Hour
}

type ConfigForPermission struct {
	EcoPkgSig                     string `json:"ecopkg_sig"`
	TCSig                         string `json:"tc_sig"`
	MinNumApprovedByTC            int    `json:"min_num_approved_by_tc"`
	MinNumApprovedBySigMaintainer int    `json:"min_num_approved_by_sig_maintainer"`
}

func (cfg *ConfigForPermission) setDefault() {
	if cfg.EcoPkgSig == "" {
		cfg.EcoPkgSig = "ecopkg"
	}

	if cfg.TCSig == "" {
		cfg.TCSig = "TC"
	}

	if cfg.MinNumApprovedByTC <= 0 {
		cfg.MinNumApprovedByTC = 1
	}

	if cfg.MinNumApprovedBySigMaintainer <= 0 {
		cfg.MinNumApprovedBySigMaintainer = 2
	}
}
