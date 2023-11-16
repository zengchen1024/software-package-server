package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	SendPkgAppliedEvent(EventMessage) error
	SendPkgCodeUpdatedEvent(EventMessage) error
	SendPkgRetestedEvent(EventMessage) error
	NotifyPkgApproved(EventMessage) error
	NotifyPkgRejected(EventMessage) error
	NotifyPkgAbandoned(EventMessage) error
	NotifyPkgAlreadyExisted(EventMessage) error
}

type SoftwarePkgIndirectMessage interface {
	SendSoftwarePkgCodeChangedEvent(EventMessage) error
}
