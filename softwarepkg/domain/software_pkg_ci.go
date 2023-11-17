package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

type pkgCI interface {
	// return new pr num
	StartNewCI(pkg *SoftwarePkg) (int, error)
}

// SoftwarePkgCI
type SoftwarePkgCI struct {
	Id        int // The Id of running CI
	status    dp.PackageCIStatus
	StartTime int64 // deal with the case that the ci is timeout
}

func (ci *SoftwarePkgCI) start(pkg *SoftwarePkg) error {
	if !ci.status.IsCIWaiting() {
		return errors.New("can't do this")
	}

	ci.status = dp.PackageCIStatusRunning
	ci.StartTime = utils.Now()

	v, err := ciInstance.StartNewCI(pkg)
	if err == nil {
		ci.Id = v
	}

	return err
}

func (ci *SoftwarePkgCI) retest() error {
	s := ci.Status()

	if s.IsCIRunning() {
		return allerror.New(allerror.ErrorCodeCIIsRunning, "ci is running")
	}

	if s.IsCIWaiting() {
		return allerror.New(allerror.ErrorCodeCIIsWaiting, "duplicate operation")
	}

	ci.status = dp.PackageCIStatusWaiting

	return nil
}

func (ci *SoftwarePkgCI) done(ciId int, success bool) error {
	if ci.Id != ciId || !ci.Status().IsCIRunning() {
		return errors.New("ignore the ci result")
	}

	if success {
		ci.status = dp.PackageCIStatusPassed
	} else {
		ci.status = dp.PackageCIStatusFailed
	}

	ci.StartTime = 0

	return nil
}

func (ci *SoftwarePkgCI) isSuccess() bool {
	return ci.status != nil && ci.status.IsCIPassed()
}

func (ci *SoftwarePkgCI) Status() dp.PackageCIStatus {
	if ci.status == nil {
		return nil
	}

	if ci.status.IsCIRunning() && ci.StartTime+timeoutOfCI < utils.Now() {
		return dp.PackageCIStatusTimeout
	}

	return ci.status
}
