package domain

import "time"

type TraceID string

type Trace struct {
	id        TraceID
	parentID  TraceID
	estimated time.Duration
	actual    time.Duration
	startedAt time.Time
	depth     uint64 // LCA(최소 공통 조상) 용도
}

func NewTrace(
	id TraceID,
	parentID TraceID,
	estimated time.Duration,
	actual time.Duration,
	startedAt time.Time,
	depth uint64,
) *Trace {
	if id == "" {
		panic("id is required")
	}

	if depth == 0 {
		panic("depth is required")
	}

	return &Trace{
		id:        id,
		parentID:  parentID,
		estimated: estimated,
		actual:    actual,
		startedAt: startedAt,
		depth:     depth,
	}
}

func (t *Trace) Clone() *Trace {
	return NewTrace(t.id, t.parentID, t.estimated, t.actual, t.startedAt, t.depth)
}

func (t *Trace) ID() TraceID {
	return t.id
}

func (t *Trace) Estimated() time.Duration {
	return t.estimated
}

func (t *Trace) Actual() time.Duration {
	return t.actual
}

func (t *Trace) ParentID() TraceID {
	return t.parentID
}

func (t *Trace) Depth() uint64 {
	return t.depth
}

func (t *Trace) SetEstimated(estimated time.Duration) *Trace {
	ret := t.Clone()
	ret.estimated = estimated

	return ret
}

func (t *Trace) SetActual(actual time.Duration) *Trace {
	ret := t.Clone()
	ret.actual = actual

	return ret
}

func (t *Trace) SetParentID(parentID TraceID) *Trace {
	ret := t.Clone()
	ret.parentID = parentID

	return ret
}
