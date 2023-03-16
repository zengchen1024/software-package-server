package main

import (
	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

type producer struct {
	topics TopicsToNotify
}

func (p *producer) NotifyPkgAlreadyClosed(e message.EventMessage) error {
	return send(p.topics.AlreadyClosedSoftwarePkg, e)
}

func (p *producer) NotifyPkgIndirectlyApproved(e message.EventMessage) error {
	return send(p.topics.IndirectlyApprovedSoftwarePkg, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, body)
}
