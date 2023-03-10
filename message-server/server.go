package main

import (
	"encoding/json"

	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type server struct {
	service app.SoftwarePkgMessageService
}

func (s *server) HandleCIChecking(data []byte) error {
	msg := new(msgToHandleCIChecking)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandleCIChecking(cmd)
}
