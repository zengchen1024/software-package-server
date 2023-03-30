package maintainerimpl

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/server-common-lib/utils"
)

type agent struct {
	sig    *sigData
	mut    sync.RWMutex
	t      utils.Timer
	loader sigLoader
}

func (instance *agent) load(link string) error {
	v, err := instance.loader.load(link, instance.sig)
	if err != nil || v == nil {
		return err
	}

	instance.mut.Lock()
	instance.sig = v
	instance.mut.Unlock()

	return nil
}

func (instance *agent) getSigData() *sigData {
	instance.mut.RLock()
	v := instance.sig // copy the pointer
	instance.mut.RUnlock()

	return v
}

func (instance *agent) start(link string, interval time.Duration) error {
	if err := instance.load(link); err != nil {
		return err
	}

	instance.t.Start(
		func() {
			if err := instance.load(link); err != nil {
				logrus.Errorf("load failed, err:%s", err.Error())
			}
		},
		interval,
		0,
	)

	return nil
}

func (instance *agent) stop() {
	instance.t.Stop()
}
