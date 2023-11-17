package dp

const (
	pkgModificationCategorySig       = "sig"
	pkgModificationCategoryCode      = "code"
	pkgModificationCategoryPkgName   = "pkg_name"
	pkgModificationCategoryPkgDesc   = "pkg_desc"
	pkgModificationCategoryUpstream  = "upstream"
	pkgModificationCategoryPkgReason = "pkg_reason"
	pkgModificationCategoryCommitter = "committer"
)

var (
	PkgModificationCategorySig       = pkgModificationCategory(pkgModificationCategorySig)
	PkgModificationCategoryCode      = pkgModificationCategory(pkgModificationCategoryCode)
	PkgModificationCategoryPkgName   = pkgModificationCategory(pkgModificationCategoryPkgName)
	PkgModificationCategoryPkgDesc   = pkgModificationCategory(pkgModificationCategoryPkgDesc)
	PkgModificationCategoryUpstream  = pkgModificationCategory(pkgModificationCategoryUpstream)
	PkgModificationCategoryPkgReason = pkgModificationCategory(pkgModificationCategoryPkgReason)
	PkgModificationCategoryCommitter = pkgModificationCategory(pkgModificationCategoryCommitter)
)

type PkgModificationCategory interface {
	PkgModificationCategory() string
}

type pkgModificationCategory string

func (v pkgModificationCategory) PkgModificationCategory() string {
	return string(v)
}
