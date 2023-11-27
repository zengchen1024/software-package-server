package main

import (
	"flag"
	"os"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/mongodb-lib/mongodblib"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/common/controller/middleware"
	"github.com/opensourceways/software-package-server/common/infrastructure/postgresql"
	"github.com/opensourceways/software-package-server/config"
	"github.com/opensourceways/software-package-server/server"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/clavalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sensitivewordsimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/translationimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/useradapterimpl"
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
	logrusutil.ComponentInit("software-package")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// Config
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	// Sig Validator
	logrus.Debugln("sig validator")
	if err := sigvalidatorimpl.Init(&cfg.SigValidator); err != nil {
		logrus.Errorf("init sig validator failed, err:%s", err.Error())

		return
	}

	defer sigvalidatorimpl.Exit()

	// mongo
	logrus.Debugln("mongo")
	if err := mongodblib.Init(&cfg.Mongo.DB); err != nil {
		logrus.Errorf("init mongo failed, err:%s", err.Error())

		return
	}

	defer mongodblib.Close()

	// Postgresql
	logrus.Debugln("pg")
	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Errorf("init db, err:%s", err.Error())

		return
	}

	// Translation
	logrus.Debugln("Translation")
	err = translationimpl.Init(
		&cfg.Translation, cfg.SoftwarePkg.DomainPrimitive.SupportedLanguages,
	)
	if err != nil {
		logrus.Errorf("init translation err:%s", err.Error())

		return
	}

	// Sensitive words
	logrus.Debugln("Sensitive")
	if err = sensitivewordsimpl.Init(&cfg.SensitiveWords); err != nil {
		logrus.Errorf("init sensitivewords err:%s", err.Error())

		return
	}

	// Encryption
	if err = utils.InitEncryption(cfg.Encryption.EncryptionKey); err != nil {
		logrus.Errorf("init encryption failed, err:%s", err.Error())

		return
	}

	// MQ
	logrus.Debugln("mq")
	if err = kfklib.Init(&cfg.MQ.Config, log, nil, "", true); err != nil {
		logrus.Errorf("init mq, err:%s", err.Error())

		return
	}

	defer kfklib.Exit()

	// Maintainer
	logrus.Debugln("maintainer")
	if err := useradapterimpl.Init(&cfg.User); err != nil {
		logrus.Errorf("init maintainer failed, err:%s", err.Error())

		return
	}

	defer useradapterimpl.Exit()

	// Domain
	domain.Init(&cfg.SoftwarePkg.Config, useradapterimpl.UserAdapter())

	dp.Init(
		&cfg.SoftwarePkg.DomainPrimitive,
		sigvalidatorimpl.SigValidator(),
		sensitivewordsimpl.Sensitive(),
	)

	// Pkg manager depends on the domain, so it should be initialized after domain
	logrus.Debugln("pkg manager")
	if err = pkgmanagerimpl.Init(&cfg.PkgManager); err != nil {
		logrus.Errorf("init pkg manager failed, err:%s", err.Error())

		return
	}

	middleware.Init(&cfg.Middleware)

	clavalidatorimpl.Init(&cfg.CLA)

	// run
	logrus.Debugln("start server")
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}
