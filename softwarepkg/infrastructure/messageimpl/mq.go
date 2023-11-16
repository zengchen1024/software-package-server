package messageimpl

import (
	"github.com/sirupsen/logrus"

	kfklib "github.com/opensourceways/kafka-lib/agent"
)

func Init(cfg *Config, log *logrus.Entry) error {
	producerInstance = &producer{cfg.Topics}

	return kfklib.Init(&cfg.Config, log, nil, "", true)
}

func Exit() {
	kfklib.Exit()
}
