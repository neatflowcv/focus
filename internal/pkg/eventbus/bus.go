package eventbus

type Bus struct {
	TaskCreated *Broker[*TaskCreatedEvent]
	TaskDeleted *Broker[*TaskDeletedEvent]
}

func NewBus() *Bus {
	return &Bus{
		TaskCreated: NewBroker[*TaskCreatedEvent](),
		TaskDeleted: NewBroker[*TaskDeletedEvent](),
	}
}
