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
	CreatedAt   time.Time
	CompletedAt sql.NullTime
}

func (t *Task) ToDomain() *domain.Task {
	return domain.NewTask(
		domain.TaskID(t.ID),
		t.Title,
		domain.TaskStatus(t.Status),
		t.CreatedAt,
		t.CompletedAt.Time,
	)
}

func FromDomainTask(task *domain.Task, username string) *Task {
	return &Task{
		ID:          string(task.ID()),
		Username:    username,
		Title:       task.Title(),
		Status:      string(task.Status()),
		CreatedAt:   task.CreatedAt(),
		CompletedAt: toNullTime(task.CompletedAt()),
	}
}

func toNullTime(time time.Time) sql.NullTime {
	return sql.NullTime{Time: time, Valid: !time.IsZero()}
}
