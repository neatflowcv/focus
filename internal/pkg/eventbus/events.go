package eventbus

type TaskCreatedEvent struct {
	TaskID   string
	ParentID string
}

type TaskDeletedEvent struct {
	TaskID string
}

type TaskRelationUpdatedEvent struct {
	TaskID      string
	OldParentID string
	NewParentID string
	OldNextID   string
	NewNextID   string
}
