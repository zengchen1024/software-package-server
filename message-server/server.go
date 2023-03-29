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
		topics.SoftwarePkgCIChecking:     s.handlePkgCIChecking,
		topics.SoftwarePkgCIChecked:      s.handlePkgCIChecked,
		topics.SoftwarePkgCodeSaved:      s.handlePkgCodeSaved,
		topics.SoftwarePkgInitialized:    s.handlePkgInitialized,
		topics.SoftwarePkgRepoCreated:    s.handlePkgRepoCreated,
		topics.SoftwarePkgAlreadyExisted: s.handlePkgAlreadyExisted,
	}

	return kafka.Subscriber().Subscribe(cfg.GroupName, h)
}

func (s *server) handlePkgCIChecking(data []byte) error {
	cmd, err := cmdToHandlePkgCIChecking(data)
	if err != nil {
		return err
	}

	return s.service.HandlePkgCIChecking(cmd)
}

func (s *server) handlePkgCIChecked(data []byte) error {
	msg := new(msgToHandlePkgCIChecked)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgCIChecked(cmd)
}

func (s *server) handlePkgInitialized(data []byte) error {
	msg := new(msgToHandlePkgInitialized)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgInitialized(cmd)
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

func (s *server) handlePkgCodeSaved(data []byte) error {
	msg := new(msgToHandlePkgCodeSaved)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgCodeSaved(cmd)
}

func (s *server) handlePkgAlreadyExisted(data []byte) error {
	msg := new(msgToHandlePkgAlreadyExisted)
	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgAlreadyExisted(cmd)
}
