package maintainerimpl

import "time"

type Config struct {
	ReadURL string `json:"read_url" required:"true"`

	// Interval the unit is hour
	Interval int `json:"interval" required:"true"`
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Duration(cfg.Interval) * time.Hour
}
