package domain

import "time"

type TraceID string

type Trace struct {
	id        TraceID
	parentID  TraceID
	estimated time.Duration
	actual    time.Duration
	startedAt time.Time
}

func NewTrace(
	id TraceID,
	parentID TraceID,
	estimated time.Duration,
	actual time.Duration,
	startedAt time.Time,
) *Trace {
	if id == "" {
		panic("id is required")
	}

	return &Trace{
		id:        id,
		parentID:  parentID,
		estimated: estimated,
		actual:    actual,
		startedAt: startedAt,
	}
}

func (t *Trace) Clone() *Trace {
	return NewTrace(t.id, t.parentID, t.estimated, t.actual, t.startedAt)
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
