package dp

const (
	pkgReviewResultRejected = "rejected"
	pkgReviewResultApproved = "approved"
)

var (
	PkgReviewResultRejected = packageReviewResult(pkgReviewResultRejected)
	PkgReviewResultApproved = packageReviewResult(pkgReviewResultApproved)
)

type PackageReviewResult interface {
	PackageReviewResult() string
}

type packageReviewResult string

func (v packageReviewResult) PackageReviewResult() string {
	return string(v)
}

func IsReviewResultRejected(r PackageReviewResult) bool {
	return r != nil && r.PackageReviewResult() == pkgReviewResultRejected
}
