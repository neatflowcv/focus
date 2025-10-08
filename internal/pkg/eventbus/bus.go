package eventbus

type Bus struct {
	TaskCreated     *Broker[*TaskCreatedEvent]
	TaskDeleted     *Broker[*TaskDeletedEvent]
	RelationCreated *Broker[*RelationCreatedEvent]
	RelationDeleted *Broker[*RelationDeletedEvent]
	RelationUpdated *Broker[*RelationUpdatedEvent]
}

func NewBus() *Bus {
	return &Bus{
		TaskCreated:     NewBroker[*TaskCreatedEvent](),
		TaskDeleted:     NewBroker[*TaskDeletedEvent](),
		RelationCreated: NewBroker[*RelationCreatedEvent](),
		RelationDeleted: NewBroker[*RelationDeletedEvent](),
		RelationUpdated: NewBroker[*RelationUpdatedEvent](),
	}
}
