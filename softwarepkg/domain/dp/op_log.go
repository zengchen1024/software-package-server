package dp

const (
	packageOpreationLogActionUpdate  = "update"
	packageOpreationLogActionReject  = "reject"
	packageOpreationLogActionApprove = "approve"
	packageOpreationLogActionAbandon = "abandon"
	packageOpreationLogActionRerunci = "rerunci"
)

var (
	PackageOpreationLogActionUpdate  = packageOpreationLogAction(packageOpreationLogActionUpdate)
	PackageOpreationLogActionReject  = packageOpreationLogAction(packageOpreationLogActionReject)
	PackageOpreationLogActionApprove = packageOpreationLogAction(packageOpreationLogActionApprove)
	PackageOpreationLogActionAbandon = packageOpreationLogAction(packageOpreationLogActionAbandon)
	PackageOpreationLogActionResunci = packageOpreationLogAction(packageOpreationLogActionRerunci)
)

type PackageOpreationLogAction interface {
	PackageOpreationLogAction() string
}

type packageOpreationLogAction string

func (p packageOpreationLogAction) PackageOpreationLogAction() string {
	return string(p)
}
