package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	NotifyPkgApplied(EventMessage) error
	NotifyPkgApproved(EventMessage) error
	NotifyPkgRejected(EventMessage) error
	NotifyPkgAbandoned(EventMessage) error
}

type SoftwarePkgIndirectMessage interface {
	NotifyPkgAlreadyClosed(EventMessage) error
	NotifyPkgIndirectlyApproved(EventMessage) error
}
