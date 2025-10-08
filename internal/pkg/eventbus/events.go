package eventbus

type TaskCreatedEvent struct {
	TaskID string
}

type TaskDeletedEvent struct {
	TaskID string
}
