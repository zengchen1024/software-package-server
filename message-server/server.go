package main

import (
	"encoding/json"

	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type server struct {
	service app.SoftwarePkgMessageService
}

func (s *server) handlePkgPRCIChecked(data []byte) error {
	msg := new(msgToHandlePkgPRCIChecked)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgPRCIChecked(cmd)
}

func (s *server) handleRepoCreated(data []byte) error {
	msg := new(msgToHandlePkgRepoCreated)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgRepoCreated(cmd)
}

func (s *server) handlePkgPRClosed(data []byte) error {
	msg := new(msgToHandlePkgPRClosed)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	return s.service.HandlePkgPRClosed(msg.toCmd())
}

func (s *server) handlePkgPRMerged(data []byte) error {
	msg := new(msgToHandlePkgPRMerged)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	return s.service.HandlePkgPRMerged(msg.toCmd())
}
