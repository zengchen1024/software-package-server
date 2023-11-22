package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgciimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/softwarepkgadapter"
	"github.com/opensourceways/software-package-server/utils"
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
	logrusutil.ComponentInit("software-package-server")
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
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	// Sig Validator
	if err = sigvalidatorimpl.Init(&cfg.SigValidator); err != nil {
		logrus.Errorf("init sig validator failed, err:%s", err.Error())

		return
	}

	defer sigvalidatorimpl.Exit()

	// Domain
	domain.Init(&cfg.SoftwarePkg.Config, nil, pkgciimpl.PkgCI())
	dp.Init(&cfg.SoftwarePkg.DomainPrimitive, sigvalidatorimpl.SigValidator(), nil)

	// Pkg manager
	if err = pkgmanagerimpl.Init(&cfg.PkgManager); err != nil {
		logrus.Errorf("init pkg manager failed, err:%s", err.Error())

		return
	}

	// mongo
	if err := mongdblib.Init(&cfg.Mongo.DB); err != nil {
		logrus.Errorf("init mongo failed, err:%s", err.Error())

		return
	}

	defer mongdblib.Close()

	// Postgresql
	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Errorf("init db, err:%s", err.Error())

		return
	}

	// Encryption
	if err = utils.InitEncryption(cfg.Encryption.EncryptionKey); err != nil {
		logrus.Errorf("init encryption failed, err:%s", err.Error())

		return
	}

	// ci
	if err = pkgciimpl.Init(&cfg.PkgCI); err != nil {
		logrus.Errorf("init pkg ci failed, err:%s", err.Error())

		return
	}

	// mq
	if err = kfklib.Init(&cfg.Kafka, log, nil, "", true); err != nil {
		logrus.Errorf("initialize mq failed, err:%v", err)

		return
	}

	defer kfklib.Exit()

	// service
	messageService := app.NewSoftwarePkgMessageService(
		nil, // TODO
		softwarepkgadapter.NewsoftwarePkgAdapter(
			mongdblib.DAO(cfg.Mongo.Collections.SoftwarePkg),
		),
		pkgmanagerimpl.Instance(),
		&producer{cfg.Topics.SoftwarePkgCodeChanged},
		repositoryimpl.NewSoftwarePkgComment(&cfg.Postgresql.Config),
	)

	// run
	run(&server{messageService}, cfg)
}

func run(s *server, cfg *Config) {
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

	if err := s.run(ctx, cfg); err != nil {
		logrus.Errorf("server exited, err:%s", err.Error())
	}
}
