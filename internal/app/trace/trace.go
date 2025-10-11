package trace

import (
	"time"
)

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

type ListTracesInput struct {
	IDs []string
}

type Trace struct {
	Estimated time.Duration
	Actual    time.Duration
	StartedAt time.Time
}

type ListTracesOutput struct {
	Traces []*Trace
}

type UpdateStatusInput struct {
	ID     string
	Status string
	Now    time.Time
}
