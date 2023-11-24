package messageimpl

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

func Producer(topic *Topics) *producer {
	return &producer{*topic}
}

type producer struct {
	topics Topics
}

func (p *producer) SendPkgAppliedEvent(e message.EventMessage) error {
	return send(p.topics.SoftwarePkgApplied, e)
}

func (p *producer) SendPkgCodeUpdatedEvent(e message.EventMessage) error {
	return send(p.topics.SoftwarePkgCodeUpdated, e)
}

func (p *producer) SendPkgRetestedEvent(e message.EventMessage) error {
	return send(p.topics.SoftwarePkgRetested, e)
}

func (p *producer) NotifyPkgApproved(e message.EventMessage) error {
	//return send(p.topics.ApprovedSoftwarePkg, e)
	return nil
}

func (p *producer) NotifyPkgRejected(e message.EventMessage) error {
	//return send(p.topics.RejectedSoftwarePkg, e)
	return nil
}

func (p *producer) NotifyPkgAbandoned(e message.EventMessage) error {
	//return send(p.topics.AbandonedSoftwarePkg, e)
	return nil
}

func (p *producer) SendPkgAlreadyExistedEvent(e message.EventMessage) error {
	return send(p.topics.SoftwarePkgAlreadyExisted, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, nil, body)
}
