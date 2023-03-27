package main

import (
	"flag"
	"os"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/config"
	"github.com/opensourceways/software-package-server/server"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sensitivewordsimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/translationimpl"
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
	logrusutil.ComponentInit("software-package")
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

	// Config
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// Postgresql
	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db, err:%s", err.Error())
	}

	// Translation
	err = translationimpl.Init(&cfg.Translation, cfg.SoftwarePkg.SupportedLanguages)
	if err != nil {
		logrus.Fatalf("init translation err:%s", err.Error())
	}

	if err = sensitivewordsimpl.Init(&cfg.SensitiveWords); err != nil {
		logrus.Fatalf("init sensitivewords err:%s", err.Error())
	}

	// MQ
	if err = messageimpl.Init(&cfg.MQ, log); err != nil {
		logrus.Fatalf("init mq, err:%s", err.Error())
	}

	defer messageimpl.Exit()

	// Sig Validator
	if err := sigvalidatorimpl.Init(&cfg.SigValidator); err != nil {
		logrus.Errorf("init sig validator failed, err:%s", err.Error())

		return
	}

	defer sigvalidatorimpl.Exit()

	// Domain
	dp.Init(&cfg.SoftwarePkg, sigvalidatorimpl.SigValidator())

	middleware.Init(&cfg.Middleware)

	pkgmanagerimpl.Init(&cfg.PkgManager)

	// run
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}
