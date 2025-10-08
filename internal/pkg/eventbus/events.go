package eventbus

type TaskCreatedEvent struct {
	TaskID string
}

type TaskDeletedEvent struct {
	TaskID string
}

type RelationCreatedEvent struct {
	RelationID string
	ParentID   string
}

type RelationDeletedEvent struct {
	RelationID string
	ParentID   string
}

type RelationUpdatedEvent struct {
	RelationID  string
	OldParentID string
	NewParentID string
}
