package domain

import "time"

type ExtraID string

type Extra struct {
	id            ExtraID
	estimatedTime time.Duration
	actualTime    time.Duration
	startedAt     time.Time
	leaf          bool
}

func NewExtra(
	id ExtraID,
	estimatedTime time.Duration,
	actualTime time.Duration,
	startedAt time.Time,
	leaf bool,
) *Extra {
	return &Extra{
		id:            id,
		estimatedTime: estimatedTime,
		actualTime:    actualTime,
		startedAt:     startedAt,
		leaf:          leaf,
	}
}

func (e *Extra) ID() ExtraID {
	return e.id
}

func (e *Extra) EstimatedTime() time.Duration {
	return e.estimatedTime
}

func (e *Extra) ActualTime() time.Duration {
	return e.actualTime
}

func (e *Extra) StartedAt() time.Time {
	return e.startedAt
}

func (e *Extra) Leaf() bool {
	return e.leaf
}

func (e *Extra) SetActualTime(actualTime time.Duration) *Extra {
	ret := e.clone()
	ret.actualTime = actualTime

	return ret
}

func (e *Extra) SetEstimatedTime(estimatedTime time.Duration) *Extra {
	ret := e.clone()
	ret.estimatedTime = estimatedTime

	return ret
}

func (e *Extra) SetLeaf(leaf bool) *Extra {
	ret := e.clone()
	ret.leaf = leaf

	return ret
}

func (e *Extra) clone() *Extra {
	return NewExtra(e.id, e.estimatedTime, e.actualTime, e.startedAt, e.leaf)
}
