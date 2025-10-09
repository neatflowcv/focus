package domain

import "time"

type Trace struct {
	estimated time.Duration
	actual    time.Duration
	accActual time.Duration
}

func NewTrace(estimated time.Duration, actual time.Duration, accActual time.Duration) *Trace {
	ret := &Trace{
		estimated: estimated,
		actual:    actual,
		accActual: accActual,
	}
	ret.validate()

	return ret
}

func (t *Trace) Clone() *Trace {
	return NewTrace(t.estimated, t.actual, t.accActual)
}

func (t *Trace) Estimated() time.Duration {
	return t.estimated
}

func (t *Trace) Actual() time.Duration {
	return t.actual
}

func (t *Trace) AccActual() time.Duration {
	return t.accActual
}

func (t *Trace) SetAccActual(accActual time.Duration) *Trace {
	ret := t.Clone()
	ret.accActual = accActual
	ret.validate()

	return ret
}

func (t *Trace) SetEstimated(estimated time.Duration) *Trace {
	ret := t.Clone()
	ret.estimated = estimated
	ret.validate()

	return ret
}

func (t *Trace) SetActual(actual time.Duration) *Trace {
	ret := t.Clone()
	ret.actual = actual
	ret.validate()

	return ret
}

func (t *Trace) validate() {
	if t.accActual < t.actual {
		panic("accActual is less than actual")
	}
}
