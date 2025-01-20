package watch

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

func NewWatchingImpl(
	repo repository.SoftwarePkg, service app.SoftwarePkgMessageService, cfg *Config,
) *WatchingImpl {
	return &WatchingImpl{
		repo:    repo,
		service: service,

		cfg:     *cfg,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type WatchingImpl struct {
	repo    repository.SoftwarePkg
	service app.SoftwarePkgMessageService

	cfg     Config
	stop    chan struct{}
	stopped chan struct{}
}

func (impl *WatchingImpl) Start() {
	go impl.do(impl.handle)
}

func (impl *WatchingImpl) Stop() {
	logrus.Info("stop watcher")

	close(impl.stop)

	<-impl.stopped

	logrus.Info("watcher stoped")
}

func (impl *WatchingImpl) do(handler func(func() bool)) {
	interval := impl.cfg.IntervalDuration()

	needStop := func() bool {
		select {
		case <-impl.stop:
			return true
		default:
			return false
		}
	}

	var timer *time.Timer

	defer func() {
		if timer != nil {
			timer.Stop()
		}

		close(impl.stopped)
	}()

	for {
		handler(needStop)

		// time starts.
		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}

		select {
		case <-impl.stop:
			return

		case <-timer.C:
		}
	}
}

func (impl *WatchingImpl) handle(needStop func() bool) {
	for times := 1; times > 0; times++ {
		if impl.doOnce(times, needStop) || needStop() {
			break
		}
	}
}

func (impl *WatchingImpl) doOnce(times int, needStop func() bool) bool {
	pkgs, err := impl.listPkgs(times)
	if err != nil {
		logrus.Errorf("list pkgs failed, err:%s", err.Error())

		return true
	}

	cmd := app.CmdToStartCI{AutoRetest: true}
	for i := range pkgs {
		cmd.PkgId = pkgs[i].Id

		if err := impl.service.StartCI(cmd); err != nil {
			logrus.Errorf("retest pkg:(%s) automatically failed, err:%s", cmd.PkgId, err.Error())
		}

		if needStop() {
			return true
		}
	}

	return len(pkgs) == 0
}

func (impl *WatchingImpl) listPkgs(pageNum int) (pkgs []repository.SoftwarePkgInfo, err error) {
	for i := 0; i < 3; i++ {
		pkgs, _, err = impl.repo.FindAll(&repository.OptToFindSoftwarePkgs{
			Phase:        dp.PackagePhaseReviewing,
			CountPerPage: 10,
			PageNum:      pageNum,
		})
		if err == nil {
			return
		}

		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return
}
