package dp

const (
	packageOperationLogActionUpdate  = "update"
	packageOperationLogActionReject  = "reject"
	packageOperationLogActionReview  = "review"
	packageOperationLogActionAbandon = "abandon"
	packageOperationLogActionRerunci = "rerunci"
)

var (
	PackageOperationLogActionUpdate  = packageOperationLogAction(packageOperationLogActionUpdate)
	PackageOperationLogActionReject  = packageOperationLogAction(packageOperationLogActionReject)
	PackageOperationLogActionReview  = packageOperationLogAction(packageOperationLogActionReview)
	PackageOperationLogActionAbandon = packageOperationLogAction(packageOperationLogActionAbandon)
	PackageOperationLogActionResunci = packageOperationLogAction(packageOperationLogActionRerunci)
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
