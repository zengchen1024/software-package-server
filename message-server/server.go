package main

import (
	"context"
	"encoding/json"

	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

const retryNum = 3

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

	err := kfklib.SubscribeWithStrategyOfSendBack(
		"software_pkg_download_code", s.downloadPkgCode,
		[]string{topics.SoftwarePkgApplied, topics.SoftwarePkgCodeUpdated},
	)
	if err != nil {
		return err
	}

	err = kfklib.SubscribeWithStrategyOfRetry(
		"software_pkg_start_ci", s.startCI,
		[]string{topics.SoftwarePkgCodeChanged, topics.SoftwarePkgRetested},
		retryNum,
	)
	if err != nil {
		return err
	}

	err = kfklib.SubscribeWithStrategyOfRetry(
		"software_pkg_ci_done", s.handlePkgCIDone,
		[]string{topics.SoftwarePkgCIDone}, retryNum,
	)
	if err != nil {
		return err
	}

	err = kfklib.SubscribeWithStrategyOfRetry(
		"software_pkg_repo_code_pushed", s.handlePkgRepoCodePushed,
		[]string{topics.SoftwarePkgRepoCodePushed}, retryNum,
	)
	if err != nil {
		return err
	}

	err = kfklib.SubscribeWithStrategyOfRetry(
		"software_pkg_closed", s.handlePkgClosed,
		[]string{topics.SoftwarePkgClosed}, retryNum,
	)
	if err != nil {
		return err
	}

	return kfklib.Subscribe(
		"software_pkg_import_existed", s.importPkg,
		[]string{topics.SoftwarePkgAlreadyExisted},
	)
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

func (s *server) handlePkgClosed(data []byte, m map[string]string) error {
	e, err := domain.UnmarshalToSoftwarePkgClosedEvent(data)
	if err != nil {
		return err
	}

	return s.service.HandlePkgClosed(&e)
}

func (s *server) importPkg(data []byte, m map[string]string) error {
	cmd, err := cmdToHandlePkgAlreadyExisted(data)
	if err != nil {
		return err
	}

	return s.service.ImportPkg(cmd)
}
