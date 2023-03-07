package messageimpl

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
)

func Init(cfg *Config, log *logrus.Entry) error {
	producerInstance = &producer{cfg.Topics}

	return kafka.Init(&cfg.Config, log)
}

func Exit() {
	kafka.Exit()
}
