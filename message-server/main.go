package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/common/infrastructure/kafka"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/maintainerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false, "whether to enable debug model.",
	)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// Postgresql
	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db, err:%s", err.Error())
	}

	// mq
	if err = kafka.Init(&cfg.Kafka, log); err != nil {
		logrus.Fatalf("initialize mq failed, err:%v", err)
	}

	defer kafka.Exit()

	messageService := app.NewSoftwarePkgMessageService(
		repositoryimpl.NewSoftwarePkg(&cfg.Postgresql.Config),
		&producer{topics: cfg.TopicsToNotify},
		maintainerimpl.NewMaintainerImpl(&cfg.Maintainer),
	)

	s := &server{messageService}
	if err := subscribe(s, cfg); err != nil {
		logrus.Errorf("subscribe failed, err:%v", err)

		return
	}

	dp.Init(&cfg.SoftwarePkg)

	// run
	run()
}

func subscribe(s *server, cfg *Config) error {
	topics := &cfg.Topics

	h := map[string]kafka.Handler{
		topics.SoftwarePkgPRMerged:    s.handlePkgPRMerged,
		topics.SoftwarePkgPRClosed:    s.handlePkgPRClosed,
		topics.SoftwarePkgPRCIChecked: s.handlePkgPRCIChecked,
		topics.SoftwarePkgRepoCreated: s.handlePkgRepoCreated,
	}

	return kafka.Subscriber().Subscribe(cfg.GroupName, h)
}

func run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")

			return

		case <-sig:
			logrus.Info("receive exit signal")
			done()
			called = true

			return
		}
	}(ctx)

	<-ctx.Done()
}
