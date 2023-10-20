package dp

const (
	checkItemCategorySig       = "sig"
	checkItemCategoryCode      = "code"
	checkItemCategoryPkgName   = "pkg_name"
	checkItemCategoryPkgDesc   = "pkg_desc"
	checkItemCategoryUpstream  = "upstream"
	checkItemCategoryPkgReason = "pkg_reason"
	checkItemCategoryCommitter = "committer"
)

var (
	CheckItemCategorySig       = checkItemCategory(checkItemCategorySig)
	CheckItemCategoryCode      = checkItemCategory(checkItemCategoryCode)
	CheckItemCategoryPkgName   = checkItemCategory(checkItemCategoryPkgName)
	CheckItemCategoryPkgDesc   = checkItemCategory(checkItemCategoryPkgDesc)
	CheckItemCategoryUpstream  = checkItemCategory(checkItemCategoryUpstream)
	CheckItemCategoryPkgReason = checkItemCategory(checkItemCategoryPkgReason)
	CheckItemCategoryCommitter = checkItemCategory(checkItemCategoryCommitter)
)

type CheckItemCategory interface {
	CheckItemCategory() string
}

type checkItemCategory string

func (v checkItemCategory) CheckItemCategory() string {
	return string(v)
}

func IsSameCheckItemCategory(c1, c2 CheckItemCategory) bool {
	return c1 != nil && c2 != nil && c1.CheckItemCategory() == c2.CheckItemCategory()
}
