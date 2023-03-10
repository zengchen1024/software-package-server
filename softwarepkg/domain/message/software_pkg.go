package message

type EventMessage interface {
	ToMessage() ([]byte, error)
}

type SoftwarePkgMessage interface {
	NotifyPkgApplied(EventMessage) error
	NotifyPkgApproved(EventMessage) error
	NotifyPkgRejected(EventMessage) error
	NotifyPkgAbandoned(EventMessage) error
	NotifyPkgAlreadyClosed(EventMessage) error
}
