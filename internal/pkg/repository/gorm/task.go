package gorm

import (
	"database/sql"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Task struct {
	ID            string `gorm:"primaryKey"`
	Username      string
	Title         string
	Status        string
	Order         float64
	CreatedAt     time.Time
	CompletedAt   sql.NullTime
	StartedAt     sql.NullTime
	EstimatedTime sql.NullInt64
	ActualTime    sql.NullInt64
}

func (t *Task) ToDomain(parentID domain.TaskID) *domain.Task {
	return domain.NewTask(
		domain.TaskID(t.ID),
		parentID,
		t.Title,
		t.CreatedAt,
		domain.TaskStatus(t.Status),
		t.Order,
		t.CompletedAt.Time,
		t.StartedAt.Time,
		time.Duration(t.EstimatedTime.Int64)*time.Second,
		time.Duration(t.ActualTime.Int64)*time.Second,
	)
}

func FromDomainTask(domainTask *domain.Task, username string) *Task {
	return &Task{
		ID:            string(domainTask.ID()),
		Username:      username,
		Title:         domainTask.Title(),
		Status:        string(domainTask.Status()),
		Order:         domainTask.Order(),
		CreatedAt:     domainTask.CreatedAt(),
		CompletedAt:   toNullTime(domainTask.CompletedAt()),
		StartedAt:     toNullTime(domainTask.StartedAt()),
		EstimatedTime: toNullInt64(domainTask.EstimatedTime()),
		ActualTime:    toNullInt64(domainTask.ActualTime()),
	}
}

func toNullTime(time time.Time) sql.NullTime {
	return sql.NullTime{Time: time, Valid: !time.IsZero()}
}

func toNullInt64(duration time.Duration) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(duration.Seconds()), Valid: duration.Seconds() != 0}
}
