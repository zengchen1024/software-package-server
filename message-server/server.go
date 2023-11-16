package main

import (
	"context"
	"encoding/json"

	kfklib "github.com/opensourceways/kafka-lib/agent"

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

	err := kfklib.Subscribe(
		"software_pkg_download_code", s.downloadPkgCode,
		[]string{topics.SoftwarePkgApplied, topics.SoftwarePkgCodeUpdated},
	)
	if err != nil {
		return err
	}

	err = kfklib.Subscribe(
		"software_pkg_start_ci", s.startCI,
		[]string{topics.SoftwarePkgCodeChanged, topics.SoftwarePkgRetested},
	)
	if err != nil {
		return err
	}

	err = kfklib.Subscribe(
		"software_pkg_ci_done", s.handlePkgCIDone,
		[]string{topics.SoftwarePkgCIDone},
	)
	if err != nil {
		return err
	}

	err = kfklib.Subscribe(
		"software_pkg_import_existed", s.importPkg,
		[]string{topics.SoftwarePkgAlreadyExisted},
	)
	if err != nil {
		return err
	}

	err = kfklib.Subscribe(
		"software_pkg_repo_code_pushed", s.handlePkgRepoCodePushed,
		[]string{topics.SoftwarePkgRepoCodePushed},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) downloadPkgCode(data []byte, m map[string]string) error {
	cmd, err := cmdToDownloadPkgCode(data)
	if err != nil {
		return err
	}

	return s.service.DownloadPkgCode(cmd)
}

func (s *server) startCI(data []byte, m map[string]string) error {
	cmd, err := cmdToStartCI(data)
	if err != nil {
		return err
	}

	return s.service.StartCI(cmd)
}

func (s *server) handlePkgCIDone(data []byte, m map[string]string) error {
	msg := new(msgToHandlePkgCIDone)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	return s.service.HandlePkgCIDone(msg.toCmd())
}

func (s *server) handlePkgRepoCodePushed(data []byte, m map[string]string) error {
	msg := new(msgToHandlePkgRepoCodePushed)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	cmd, err := msg.toCmd()
	if err != nil {
		return err
	}

	return s.service.HandlePkgRepoCodePushed(cmd)
}

func (s *server) importPkg(data []byte, m map[string]string) error {
	cmd, err := cmdToHandlePkgAlreadyExisted(data)
	if err != nil {
		return err
	}

	return s.service.ImportPkg(cmd)
}
