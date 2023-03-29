package pkgmanagerimpl

type Config struct {
	Org         string     `json:"org"            required:"true"`
	Base        BaseConfig `json:"base"           required:"true"`
	AccessToken string     `json:"access_token"   required:"true"`
}

type BaseConfig struct {
	Org    string `json:"org"    required:"true"`
	Repo   string `json:"repo"   required:"true"`
	Branch string `json:"branch" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.Org == "" {
		cfg.Org = "src-openeuler"
	}
}

func (cfg *Config) Token() func() []byte {
	return func() []byte {
		return []byte(cfg.AccessToken)
	}
}
