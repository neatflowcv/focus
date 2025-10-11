package eventbus

type Bus struct {
	TaskCreated         *Broker[*TaskCreatedEvent]
	TaskDeleted         *Broker[*TaskDeletedEvent]
	TaskRelationUpdated *Broker[*TaskRelationUpdatedEvent]
}

func NewBus() *Bus {
	return &Bus{
		TaskCreated:         NewBroker[*TaskCreatedEvent](),
		TaskDeleted:         NewBroker[*TaskDeletedEvent](),
		TaskRelationUpdated: NewBroker[*TaskRelationUpdatedEvent](),
	}
}
