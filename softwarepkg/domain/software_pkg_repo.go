package domain

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Platform   dp.PackagePlatform
	Committers []dp.Account
}

func (r *SoftwarePkgRepo) isCommitter(u dp.Account) bool {
	for i := range r.Committers {
		if dp.IsSameAccount(r.Committers[i], u) {
			return true
		}
	}

	return false
}

func (r *SoftwarePkgRepo) update(r1 *SoftwarePkgRepo) string {
	if p := r1.Platform; p != nil && r.Platform.PackagePlatform() != p.PackagePlatform() {
		r.Platform = p
	}

	if len(r1.Committers) > 0 && !r.isSameCommitters(r1) {
		r.Committers = r1.Committers

		return pkgModificationCommitter
	}

	return ""
}

func (r *SoftwarePkgRepo) isSameCommitters(r1 *SoftwarePkgRepo) bool {
	if len(r.Committers) != len(r1.Committers) {
		return false
	}

	m := r.committersMap()

	for i := range r1.Committers {
		if !m[r1.Committers[i].Account()] {
			return false
		}
	}

	return true
}

func (r *SoftwarePkgRepo) committersMap() map[string]bool {
	m := map[string]bool{}

	for i := range r.Committers {
		m[r.Committers[i].Account()] = true
	}

	return m
}
