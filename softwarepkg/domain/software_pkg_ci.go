package domain

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

type pkgCI interface {
	// return new ci id(pr num)
	StartNewCI(pkg *SoftwarePkg) (int, error)

	// should return nil if the CI had been cleared(pr is closed)
	ClearCI(int) error
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

func (ci *SoftwarePkgCI) init() {
	ci.Id = 0
	ci.status = dp.PackageCIStatusWaiting
	ci.StartTime = utils.Now()
}

func (ci *SoftwarePkgCI) reset() error {
	if ci.Id != 0 && !(ci.status.IsCIFailed() || ci.status.IsCIPassed()) {
		if err := ciInstance.ClearCI(ci.Id); err != nil {
			return err
		}
	}

	ci.init()

	return nil
}

func (ci *SoftwarePkgCI) start(pkg *SoftwarePkg) error {
	if !ci.status.IsCIWaiting() {
		return errors.New("ci is not waiting")
	}

	return ci.startCI(pkg)
}

func (ci *SoftwarePkgCI) startCI(pkg *SoftwarePkg) error {
	v, err := ciInstance.StartNewCI(pkg)
	if err == nil {
		ci.Id = v
		ci.status = dp.PackageCIStatusRunning
		ci.StartTime = utils.Now()
	}

	return err
}

func (ci *SoftwarePkgCI) retest(pkg *SoftwarePkg) error {
	s := ci.Status()

	if s.IsCIRunning() {
		return allerror.New(allerror.ErrorCodeCIIsRunning, "ci is running")
	}

	// no need to test if it is passed.
	if s.IsCIPassed() {
		return allerror.New(allerror.ErrorCodeCIIsPassed, "ci is passed")
	}

	if ci.Id != 0 && !s.IsCIFailed() {
		if err := ciInstance.ClearCI(ci.Id); err != nil {
			return err
		}
	}

	// It is ok to start ci here even if retesting cocurrently,
	// because there is only one pr to be created successfully at last.
	return ci.startCI(pkg)
}

func (ci *SoftwarePkgCI) autoRetest(pkg *SoftwarePkg) error {
	s := ci.Status()

	if s.IsCIRunning() || s.IsCIPassed() || s.IsCIFailed() {
		return errors.New("can't retest automatically, because the status is not correct")
	}

	// ci will be triggered by event, so wait until time is up.
	if s.IsCIWaiting() && utils.Now() < ci.StartTime+ciConfig.CIWaitTimeout {
		return errors.New("can't retest automatically, because it is not time yet")
	}

	if ci.Id != 0 {
		if err := ciInstance.ClearCI(ci.Id); err != nil {
			return err
		}
	}

	return ci.startCI(pkg)
}

func (ci *SoftwarePkgCI) done(ciId int, success bool) error {
	if err := ciInstance.ClearCI(ciId); err != nil {
		return err
	}

	if ci.Id != ciId || !ci.status.IsCIRunning() {
		logrus.Errorf(
			"ignore the ci result, ci id %v != %v, status=%v",
			ci.Id, ciId, ci.status,
		)

		return allerror.New(allerror.ErrorCodeCIIsUnmatched, "ignore the ci result")
	}

	if success {
		ci.status = dp.PackageCIStatusPassed
	} else {
		ci.status = dp.PackageCIStatusFailed
	}

	ci.StartTime = 0

	return nil
}

func (ci *SoftwarePkgCI) isPassed() bool {
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
