package gorm

import (
	"database/sql"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Trace struct {
	ID        string
	ParentID  sql.NullString
	Estimated sql.NullInt64
	Actual    sql.NullInt64
	StartedAt sql.NullTime
}

func FromDomainTrace(trace *domain.Trace) *Trace {
	return &Trace{
		ID:        string(trace.ID()),
		ParentID:  sql.NullString{String: string(trace.ParentID()), Valid: true},
		Estimated: sql.NullInt64{Int64: int64(trace.Estimated()), Valid: true},
		Actual:    sql.NullInt64{Int64: int64(trace.Actual()), Valid: true},
		StartedAt: sql.NullTime{Time: trace.StartedAt(), Valid: true},
	}
}

func (t *Trace) ToDomain() *domain.Trace {
	return domain.NewTrace(
		domain.TraceID(t.ID),
		domain.TraceID(getString(t.ParentID)),
		time.Duration(getInt64(t.Estimated)),
		time.Duration(getInt64(t.Actual)),
		getTime(t.StartedAt),
	)
}

func getTime(s sql.NullTime) time.Time {
	if !s.Valid {
		return time.Time{}
	}

	return s.Time
}

func getInt64(s sql.NullInt64) int64 {
	if !s.Valid {
		return 0
	}

	return s.Int64
}
