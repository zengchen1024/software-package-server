package dp

const (
	packageOperationLogActionUpdate  = "update"
	packageOperationLogActionReject  = "reject"
	packageOperationLogActionReview  = "review"
	packageOperationLogActionRetest  = "retest"
	packageOperationLogActionAbandon = "abandon"
)

var (
	PackageOperationLogActionUpdate  = packageOperationLogAction(packageOperationLogActionUpdate)
	PackageOperationLogActionReject  = packageOperationLogAction(packageOperationLogActionReject)
	PackageOperationLogActionReview  = packageOperationLogAction(packageOperationLogActionReview)
	PackageOperationLogActionRetest  = packageOperationLogAction(packageOperationLogActionRetest)
	PackageOperationLogActionAbandon = packageOperationLogAction(packageOperationLogActionAbandon)
)

type PackageOperationLogAction interface {
	PackageOperationLogAction() string
}

func NewPackageOperationLogAction(action string) PackageOperationLogAction {
	return packageOperationLogAction(action)
}

type packageOperationLogAction string

func (p packageOperationLogAction) PackageOperationLogAction() string {
	return string(p)
}
