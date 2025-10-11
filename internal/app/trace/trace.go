package trace

import "time"

type CreateTraceInput struct {
	ID       string
	ParentID string
}

type DeleteTraceInput struct {
	ID string
}

type SetActualInput struct {
	ID     string
	Actual time.Duration
}

type UpdateParentInput struct {
	ID       string
	ParentID string
}
