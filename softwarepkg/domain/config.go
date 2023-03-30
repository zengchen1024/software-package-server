package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MinNumOfApprovers int `json:"min_num_of_approvers"`
}

func (cfg *Config) SetDefault() {
	if cfg.MinNumOfApprovers <= 0 {
		cfg.MinNumOfApprovers = 1
	}
}
