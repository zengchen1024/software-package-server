package messageimpl

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

var producerInstance *producer

func Producer() *producer {
	return producerInstance
}

type producer struct {
	topics Topics
}

func (p *producer) NotifyPkgApplied(e *domain.SoftwarePkgAppliedEvent) error {
	return send(p.topics.ApplyingSoftwarePkg, e)
}

func (p *producer) NotifyPkgApproved(e *domain.SoftwarePkgApprovedEvent) error {
	return send(p.topics.ApprovedSoftwarePkg, e)
}

func (p *producer) NotifyPkgRejected(e *domain.SoftwarePkgRejectedEvent) error {
	return send(p.topics.RejectedSoftwarePkg, e)
}

// send
type event interface {
	ToMessage() ([]byte, error)
}

func send(topic string, v event) error {
	body, err := v.ToMessage()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, &mq.Message{
		Body: body,
	})
}
