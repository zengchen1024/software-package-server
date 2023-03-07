package messageimpl

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

var producerInstance *producer

func Producer() *producer {
	return producerInstance
}

type producer struct {
	topics Topics
}

func (p *producer) NotifyPkgApplied(e message.EventMessage) error {
	return send(p.topics.ApplyingSoftwarePkg, e)
}

func (p *producer) NotifyPkgApproved(e message.EventMessage) error {
	return send(p.topics.ApprovedSoftwarePkg, e)
}

func (p *producer) NotifyPkgRejected(e message.EventMessage) error {
	return send(p.topics.RejectedSoftwarePkg, e)
}

func (p *producer) NotifyPkgAbandoned(e message.EventMessage) error {
	return send(p.topics.AbandonedSoftwarePkg, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.ToMessage()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, &mq.Message{
		Body: body,
	})
}
