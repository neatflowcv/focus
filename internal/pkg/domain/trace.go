package domain

import "time"

type Trace struct {
	estimated time.Duration
	actual    time.Duration
}

func NewTrace(estimated time.Duration, actual time.Duration) *Trace {
	return &Trace{
		estimated: estimated,
		actual:    actual,
	}
}

func (t *Trace) Clone() *Trace {
	return NewTrace(t.estimated, t.actual)
}

func (t *Trace) Estimated() time.Duration {
	return t.estimated
}

func (t *Trace) Actual() time.Duration {
	return t.actual
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
