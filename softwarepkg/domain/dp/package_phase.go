package dp

import "errors"

const (
	packagePhaseCreatingRepo = "creating_repo"
	packagePhaseReviewing    = "reviewing"
	packagePhaseImported     = "imported"
	packagePhaseClosed       = "closed"
)

var (
	validPackagePhase = map[string]bool{
		packagePhaseCreatingRepo: true,
		packagePhaseReviewing:    true,
		packagePhaseImported:     true,
		packagePhaseClosed:       true,
	}

	PackagePhaseCreatingRepo = packagePhase(packagePhaseCreatingRepo)
	PackagePhaseReviewing    = packagePhase(packagePhaseReviewing)
	PackagePhaseImported     = packagePhase(packagePhaseImported)
	PackagePhaseClosed       = packagePhase(packagePhaseClosed)
)

type PackagePhase interface {
	PackagePhase() string
	IsClosed() bool
	IsReviewing() bool
	IsCreatingRepo() bool
}

func NewPackagePhase(v string) (PackagePhase, error) {
	if !validPackagePhase[v] {
		return nil, errors.New("invalid package phase")
	}

	return packagePhase(v), nil
}

type packagePhase string

func (v packagePhase) PackagePhase() string {
	return string(v)
}

func (v packagePhase) IsClosed() bool {
	return string(v) == packagePhaseClosed
}

func (v packagePhase) IsReviewing() bool {
	return string(v) == packagePhaseReviewing
}

func (v packagePhase) IsCreatingRepo() bool {
	return string(v) == packagePhaseCreatingRepo
}
