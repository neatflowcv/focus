package gorm

import (
	"database/sql"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Task struct {
	ID          string `gorm:"primaryKey"`
	Username    string
	Title       string
	Status      string
	Order       float64
	CreatedAt   time.Time
	CompletedAt sql.NullTime
}

func (t *Task) ToDomain(parentID domain.TaskID) *domain.Task {
	return domain.NewTask(
		domain.TaskID(t.ID),
		parentID,
		t.Title,
		domain.TaskStatus(t.Status),
		t.Order,
		t.CreatedAt,
		t.CompletedAt.Time,
	)
}

func FromDomainTask(domainTask *domain.Task, username string) *Task {
	return &Task{
		ID:          string(domainTask.ID()),
		Username:    username,
		Title:       domainTask.Title(),
		Status:      string(domainTask.Status()),
		Order:       domainTask.Order(),
		CreatedAt:   domainTask.CreatedAt(),
		CompletedAt: toNullTime(domainTask.CompletedAt()),
	}
}

func toNullTime(time time.Time) sql.NullTime {
	return sql.NullTime{Time: time, Valid: !time.IsZero()}
}
