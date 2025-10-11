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
	ret := &Trace{
		id:        id,
		parentID:  parentID,
		estimated: estimated,
		actual:    actual,
		startedAt: startedAt,
	}
	ret.validate()

	return ret
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

func (t *Trace) StartedAt() time.Time {
	return t.startedAt
}

func (t *Trace) SetEstimated(estimated time.Duration) *Trace {
	ret := t.clone()
	ret.estimated = estimated
	ret.validate()

	return ret
}

func (t *Trace) SetActual(actual time.Duration) *Trace {
	ret := t.clone()
	ret.actual = actual
	ret.validate()

	return ret
}

func (t *Trace) SetParentID(parentID TraceID) *Trace {
	ret := t.clone()
	ret.parentID = parentID
	ret.validate()

	return ret
}

func (t *Trace) validate() {
	if t.estimated < 0 {
		panic("estimated is required")
	}
}

func (t *Trace) clone() *Trace {
	return NewTrace(t.id, t.parentID, t.estimated, t.actual, t.startedAt)
}
