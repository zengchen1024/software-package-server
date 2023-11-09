package dp

import "errors"

const (
	packageCIStatusFailed  = "ci-failed"
	packageCIStatusPassed  = "ci-passed"
	packageCIStatusRunning = "ci-running"
	packageCIStatusWaiting = "ci-waiting"
	packageCIStatusTimeout = "ci-timeout"
)

var (
	validPackageCIStatus = map[string]bool{
		packageCIStatusFailed:  true,
		packageCIStatusPassed:  true,
		packageCIStatusRunning: true,
		packageCIStatusWaiting: true,
		packageCIStatusTimeout: true,
	}

	PackageCIStatusFailed  = packageCIStatus(packageCIStatusFailed)
	PackageCIStatusPassed  = packageCIStatus(packageCIStatusPassed)
	PackageCIStatusRunning = packageCIStatus(packageCIStatusRunning)
	PackageCIStatusWaiting = packageCIStatus(packageCIStatusWaiting)
	PackageCIStatusTimeout = packageCIStatus(packageCIStatusTimeout)
)

type packageCIStatus string

type PackageCIStatus interface {
	PackageCIStatus() string
	IsCIFailed() bool
	IsCIPassed() bool
	IsCIRunning() bool
	IsCIWaiting() bool
	IsCITimeout() bool
}

func NewPackageCIStatus(v string) (PackageCIStatus, error) {
	if !validPackageCIStatus[v] {
		return nil, errors.New("invalid package ci status")
	}

	return packageCIStatus(v), nil
}

func (p packageCIStatus) PackageCIStatus() string {
	return string(p)
}

func (p packageCIStatus) IsCIFailed() bool {
	return p == PackageCIStatusFailed
}

func (p packageCIStatus) IsCIPassed() bool {
	return p == PackageCIStatusPassed
}

func (p packageCIStatus) IsCIRunning() bool {
	return p == PackageCIStatusRunning
}

func (p packageCIStatus) IsCIWaiting() bool {
	return p == PackageCIStatusWaiting
}

func (p packageCIStatus) IsCITimeout() bool {
	return p == PackageCIStatusTimeout
}
