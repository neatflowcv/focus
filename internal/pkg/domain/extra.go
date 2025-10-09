package domain

import "time"

type ExtraID string

type Extra struct {
	id        ExtraID
	parentID  ExtraID
	durations *Trace
	startedAt time.Time
	leaf      bool
	status    TaskStatus
}

func NewExtra(
	id ExtraID,
	parentID ExtraID,
	durations *Trace,
	startedAt time.Time,
	leaf bool,
	status TaskStatus,
) *Extra {
	return &Extra{
		id:        id,
		parentID:  parentID,
		durations: durations,
		startedAt: startedAt,
		leaf:      leaf,
		status:    status,
	}
}

func (e *Extra) ID() ExtraID {
	return e.id
}

func (e *Extra) EstimatedTime() time.Duration {
	return e.durations.estimated
}

func (e *Extra) ActualTime() time.Duration {
	return e.durations.actual
}

func (e *Extra) StartedAt() time.Time {
	return e.startedAt
}

func (e *Extra) Leaf() bool {
	return e.leaf
}

func (e *Extra) ParentID() ExtraID {
	return e.parentID
}

func (e *Extra) Status() TaskStatus {
	return e.status
}

func (e *Extra) AccActualTime() time.Duration {
	return e.durations.accActual
}

func (e *Extra) IsCompleted() bool {
	return e.status == TaskStatusDone
}

func (e *Extra) SetActualTime(actualTime time.Duration) *Extra {
	ret := e.clone()
	ret.durations = ret.durations.SetActual(actualTime)
	ret.validate()

	return ret
}

func (e *Extra) SetEstimatedTime(estimatedTime time.Duration) *Extra {
	ret := e.clone()
	ret.durations = ret.durations.SetEstimated(estimatedTime)
	ret.validate()

	return ret
}

func (e *Extra) SetLeaf(leaf bool) *Extra {
	ret := e.clone()
	ret.leaf = leaf
	ret.validate()

	return ret
}

func (e *Extra) SetParentID(parentID ExtraID) *Extra {
	ret := e.clone()
	ret.parentID = parentID
	ret.validate()

	return ret
}

func (e *Extra) SetStatus(status TaskStatus) *Extra {
	ret := e.clone()
	ret.status = status
	ret.validate()

	return ret
}

func (e *Extra) SetAccActualTime(accActualTime time.Duration) *Extra {
	ret := e.clone()
	ret.durations = ret.durations.SetAccActual(accActualTime)
	ret.validate()

	return ret
}

func (e *Extra) clone() *Extra {
	return NewExtra(e.id, e.parentID, e.durations, e.startedAt, e.leaf, e.status)
}

func (e *Extra) validate() {
	e.status.validate()
}
