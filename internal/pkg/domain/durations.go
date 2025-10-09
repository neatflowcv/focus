package domain

import "time"

type Durations struct {
	estimated time.Duration
	actual    time.Duration
}

func NewDurations(estimated time.Duration, actual time.Duration) *Durations {
	return &Durations{
		estimated: estimated,
		actual:    actual,
	}
}

func (d *Durations) Clone() *Durations {
	return NewDurations(d.estimated, d.actual)
}

func (d *Durations) Estimated() time.Duration {
	return d.estimated
}

func (d *Durations) Actual() time.Duration {
	return d.actual
}

func (d *Durations) SetEstimated(estimated time.Duration) *Durations {
	ret := d.Clone()
	ret.estimated = estimated

	return ret
}

func (d *Durations) SetActual(actual time.Duration) *Durations {
	ret := d.Clone()
	ret.actual = actual

	return ret
}
