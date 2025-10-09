package domain

import "time"

type ExtraID string

type Extra struct {
	id        ExtraID
	parentID  ExtraID
	durations *Durations
	startedAt time.Time
	leaf      bool
}

func NewExtra(
	id ExtraID,
	parentID ExtraID,
	durations *Durations,
	startedAt time.Time,
	leaf bool,
) *Extra {
	return &Extra{
		id:        id,
		parentID:  parentID,
		durations: durations,
		startedAt: startedAt,
		leaf:      leaf,
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

func (e *Extra) SetActualTime(actualTime time.Duration) *Extra {
	ret := e.clone()
	ret.durations = ret.durations.SetActual(actualTime)

	return ret
}

func (e *Extra) SetEstimatedTime(estimatedTime time.Duration) *Extra {
	ret := e.clone()
	ret.durations = ret.durations.SetEstimated(estimatedTime)

	return ret
}

func (e *Extra) SetLeaf(leaf bool) *Extra {
	ret := e.clone()
	ret.leaf = leaf

	return ret
}

func (e *Extra) SetParentID(parentID ExtraID) *Extra {
	ret := e.clone()
	ret.parentID = parentID

	return ret
}

func (e *Extra) clone() *Extra {
	return NewExtra(e.id, e.parentID, e.durations, e.startedAt, e.leaf)
}
