package cacheagent

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/server-common-lib/utils"
)

type Loader interface {
	Load(interface{}) (interface{}, error)
}

func NewCacheAgent(loader Loader, interval time.Duration) (*Agent, error) {
	v := &Agent{
		t:      utils.NewTimer(),
		loader: loader,
	}

	err := v.start(interval)

	return v, err
}

// Agent
type Agent struct {
	// data must be a pointer
	data   interface{}
	mut    sync.RWMutex
	t      utils.Timer
	loader Loader
}

func (instance *Agent) load() error {
	v, err := instance.loader.Load(instance.data)
	if err != nil || v == nil {
		return err
	}

	instance.mut.Lock()
	instance.data = v
	instance.mut.Unlock()

	return nil
}

func (instance *Agent) GetData() interface{} {
	instance.mut.RLock()
	v := instance.data // copy the pointer
	instance.mut.RUnlock()

	return v
}

func (instance *Agent) start(interval time.Duration) error {
	if err := instance.load(); err != nil {
		return err
	}

	instance.t.Start(
		func() {
			if err := instance.load(); err != nil {
				logrus.Errorf("load failed, err:%s", err.Error())
			}
		},
		interval,
		0,
	)

	return nil
}

func (instance *Agent) Stop() {
	instance.t.Stop()
}
