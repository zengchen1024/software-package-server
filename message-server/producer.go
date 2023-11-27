package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

type producer struct {
	softwarePkgCodeChanged string
}

func (p *producer) SendPkgCodeChangedEvent(e message.EventMessage) error {
	body, err := e.Message()
	if err != nil {
		return err
	}

	return kfklib.Publish(p.softwarePkgCodeChanged, nil, body)
}
