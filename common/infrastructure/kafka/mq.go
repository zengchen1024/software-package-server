package kafka

import (
	"context"

	libkafka "github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"
)

func Init(cfg *Config, log *logrus.Entry) error {
	err := libkafka.Init(
		mq.Addresses(cfg.mqConfig().Addresses...),
		mq.Log(log),
	)
	if err != nil {
		return err
	}

	return libkafka.Connect()
}

func Exit() {
	if err := libkafka.Disconnect(); err != nil {
		logrus.Errorf("exit kafka, err:%v", err)
	}
}

func Publish(topic string, msg []byte) error {
	return libkafka.Publish(topic, &mq.Message{
		Body: msg,
	})
}

type Handler func([]byte) error

func Subscribe(ctx context.Context, handlers map[string]Handler) error {
	subscribers := make(map[string]mq.Subscriber)

	defer func() {
		for k, s := range subscribers {
			if err := s.Unsubscribe(); err != nil {
				logrus.Errorf(
					"failed to unsubscribe for topic:%s, err:%v",
					k, err,
				)
			}
		}
	}()

	// subscribe
	for topic, h := range handlers {
		s, err := registerHandler(topic, h)
		if err != nil {
			return err
		}

		if s != nil {
			subscribers[s.Topic()] = s
		} else {
			logrus.Infof("does not subscribe topic:%s", topic)
		}
	}

	// register end
	if len(subscribers) == 0 {
		return nil
	}

	<-ctx.Done()

	return nil
}

func registerHandler(topic string, h Handler) (mq.Subscriber, error) {
	if h == nil {
		return nil, nil
	}

	return libkafka.Subscribe(topic, func(e mq.Event) error {
		msg := e.Message()
		if msg == nil {
			return nil
		}

		return h(msg.Body)
	})
}
