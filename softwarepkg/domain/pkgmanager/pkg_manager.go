package pkgmanager

type PkgManager interface {
	IsPkgExisted(string) bool
}
