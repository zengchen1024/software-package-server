package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	SendPkgClosedEvent(EventMessage) error
	SendPkgAppliedEvent(EventMessage) error
	SendPkgCodeUpdatedEvent(EventMessage) error
	SendPkgRetestedEvent(EventMessage) error
	SendPkgAlreadyExistedEvent(EventMessage) error
}

type SoftwarePkgIndirectMessage interface {
	SendPkgCodeChangedEvent(EventMessage) error
}

type SoftwarePkgInitMessage interface {
	SendPkgInitializedEvent(EventMessage) error
}
