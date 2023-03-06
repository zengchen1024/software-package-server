package messageimpl

import (
	"errors"
	"regexp"
	"strings"

	"github.com/opensourceways/community-robot-lib/mq"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9][0-9]*$`)

type Config struct {
	Topics  Topics `json:"topics"  required:"true"`
	Address string `json:"address" required:"true"`
}

func (cfg *Config) mqConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.ParseAddress(),
	}
}

func (cfg *Config) Validate() error {
	if r := cfg.ParseAddress(); len(r) == 0 {
		return errors.New("invalid mq address")
	}

	return nil
}

func (cfg *Config) ParseAddress() []string {
	v := strings.Split(cfg.Address, ",")
	r := make([]string, 0, len(v))
	for i := range v {
		if reIpPort.MatchString(v[i]) {
			r = append(r, v[i])
		}
	}

	return r
}

type Topics struct {
	ApplyingSoftwarePkg  string `json:"applying_software_pkg"  required:"true"`
	ApprovedSoftwarePkg  string `json:"approved_software_pkg"  required:"true"`
	RejectedSoftwarePkg  string `json:"rejected_software_pkg"  required:"true"`
	AbandonedSoftwarePkg string `json:"rejected_software_pkg"  required:"true"`
}
