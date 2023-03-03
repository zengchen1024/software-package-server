package messageimpl

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"
)

func Init(cfg *Config, log *logrus.Entry) error {
	producerInstance = &producer{cfg.Topics}

	err := kafka.Init(
		mq.Addresses(cfg.mqConfig().Addresses...),
		mq.Log(log),
	)
	if err != nil {
		return err
	}

	return kafka.Connect()
}

func Exit(log *logrus.Entry) {
	if err := kafka.Disconnect(); err != nil {
		log.Errorf("exit kafka, err:%v", err)
	}
}
