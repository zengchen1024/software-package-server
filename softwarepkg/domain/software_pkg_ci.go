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

func NewSoftwarePkgCI(cid int, status dp.PackageCIStatus, startTime int64) SoftwarePkgCI {
	return SoftwarePkgCI{
		Id:        cid,
		status:    status,
		StartTime: startTime,
	}
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

	v, err := ciInstance.StartNewCI(pkg)
	if err == nil {
		ci.Id = v
		ci.status = dp.PackageCIStatusRunning
		ci.StartTime = utils.Now()
	}

	return err
}

func (ci *SoftwarePkgCI) retest() error {
	s := ci.Status()

	if s.IsCIRunning() {
		return allerror.New(allerror.ErrorCodeCIIsRunning, "ci is running")
	}

	if s.IsCIWaiting() {
		now := utils.Now()
		if now < ci.StartTime+ciConfig.CIWaitTimeout {
			return allerror.New(allerror.ErrorCodeCIIsWaiting, "duplicate operation")
		}
		ci.StartTime = now

		return nil
	}

	ci.status = dp.PackageCIStatusWaiting
	ci.StartTime = utils.Now()

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

	if ci.status.IsCIRunning() && ci.StartTime+ciConfig.CITimeout < utils.Now() {
		return dp.PackageCIStatusTimeout
	}

	return ci.status
}
