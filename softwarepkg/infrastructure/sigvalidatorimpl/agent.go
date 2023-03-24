package sigvalidatorimpl

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/server-common-lib/utils"
)

type agent struct {
	sig    *sigData
	md5Sum string
	mut    sync.RWMutex
	t      utils.Timer
	loader sigLoader
}

func (ca *agent) load(link string) error {
	c, md5Sum, err := ca.loader.load(link, ca.md5Sum)
	if err != nil || c == nil {
		return err
	}

	ca.mut.Lock()
	ca.sig = c
	ca.md5Sum = md5Sum
	ca.mut.Unlock()

	return nil
}

func (ca *agent) getSigData() *sigData {
	ca.mut.RLock()
	c := ca.sig // copy the pointer
	ca.mut.RUnlock()

	return c
}

func (ca *agent) start(link string, interval time.Duration) error {
	if err := ca.load(link); err != nil {
		return err
	}

	ca.t.Start(
		func() {
			if err := ca.load(link); err != nil {
				logrus.Errorf("load failed, err:%s", err.Error())
			}
		},
		interval,
		0,
	)

	return nil
}

func (ca *agent) Stop() {
	ca.t.Stop()
}
