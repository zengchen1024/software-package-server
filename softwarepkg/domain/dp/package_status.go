package dp

import "errors"

const (
	packageStatusCreatingRepo = "creating_repo"
	packageStatusInProgress   = "in-progress"
	packageStatusRejected     = "rejected"
	packageStatusGivingUp     = "giving_up"
	packageStatusNormal       = "normal"
)

var (
	validPackageStatus = map[string]bool{
		packageStatusCreatingRepo: true,
		packageStatusInProgress:   true,
		packageStatusRejected:     true,
		packageStatusGivingUp:     true,
		packageStatusNormal:       true,
	}

	PackageStatusCreatingRepo = packageStatus(packageStatusCreatingRepo)
	PackageStatusInProgress   = packageStatus(packageStatusInProgress)
	PackageStatusRejected     = packageStatus(packageStatusRejected)
	PackageStatusGivingUp     = packageStatus(packageStatusGivingUp)
	PackageStatusNormal       = packageStatus(packageStatusNormal)
)

type PackageStatus interface {
	PackageStatus() string
}

func NewPackageStatus(v string) (PackageStatus, error) {
	if !validPackageStatus[v] {
		return nil, errors.New("invalid package status")
	}

	return packageStatus(v), nil
}

type packageStatus string

func (v packageStatus) PackageStatus() string {
	return string(v)
}
