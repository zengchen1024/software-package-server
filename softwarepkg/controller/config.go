package controller

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxPageNum         int `json:"max_page_num"`
	MaxCountPerPage    int `json:"max_count_per_page"`
	MaxNumOfCommitters int `json:"max_num_of_committers"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxPageNum <= 0 {
		cfg.MaxPageNum = 10000
	}

	if cfg.MaxCountPerPage <= 0 {
		cfg.MaxCountPerPage = 100
	}

	if cfg.MaxNumOfCommitters <= 0 {
		cfg.MaxNumOfCommitters = 2
	}
}
