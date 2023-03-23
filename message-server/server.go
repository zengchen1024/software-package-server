package main

import (
	"context"
	"encoding/json"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
)

type server struct {
	service app.SoftwarePkgMessageService
}

func (s *server) run(ctx context.Context, cfg *Config) error {
	if err := s.subscribe(cfg); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *server) subscribe(cfg *Config) error {
	topics := &cfg.Topics

	h := map[string]kafka.Handler{
		topics.SoftwarePkgPRMerged:    s.handlePkgPRMerged,
		topics.SoftwarePkgPRClosed:    s.handlePkgPRClosed,
		topics.SoftwarePkgPRCIChecked: s.handlePkgPRCIChecked,
		topics.SoftwarePkgRepoCreated: s.handlePkgRepoCreated,
	}

	return kafka.Subscriber().Subscribe(cfg.GroupName, h)
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

func (s *server) handlePkgRepoCreated(data []byte) error {
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
