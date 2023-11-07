package useradapterimpl

import "time"

type Config struct {
	TCSig   string `json:"tc_sig"`
	ReadURL string `json:"read_url" required:"true"`

	// Interval the unit is hour
	Interval int `json:"interval" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.TCSig == "" {
		cfg.TCSig = "TC"
	}
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Duration(cfg.Interval) * time.Hour
}
