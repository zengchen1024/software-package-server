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
		packagePhaseClosed:       true,
		packagePhaseImported:     true,
	}

	PackagePhaseCreatingRepo = packagePhase(packagePhaseCreatingRepo)
	PackagePhaseReviewing    = packagePhase(packagePhaseReviewing)
	PackagePhaseClosed       = packagePhase(packagePhaseClosed)
	PackagePhaseImported     = packagePhase(packagePhaseImported)
)

type PackagePhase interface {
	PackagePhase() string
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

func IsPackagePhaseReviewing(s PackagePhase) bool {
	return s != nil && s.PackagePhase() == packagePhaseReviewing
}
