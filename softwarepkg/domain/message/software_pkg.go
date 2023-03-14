package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	NotifyPkgApplied(EventMessage) error
	NotifyPkgPRMerged(EventMessage) error
	NotifyPkgApproved(EventMessage) error
	NotifyPkgRejected(EventMessage) error
	NotifyPkgAbandoned(EventMessage) error
	NotifyPkgAlreadyClosed(EventMessage) error
}
