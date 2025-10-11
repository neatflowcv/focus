package gorm

import (
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Task struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	Title     string
	CreatedAt time.Time
}

func (t *Task) ToDomain() *domain.Task {
	return domain.NewTask(
		domain.TaskID(t.ID),
		t.Title,
		t.CreatedAt,
	)
}

func FromDomainTask(task *domain.Task, username string) *Task {
	return &Task{
		ID:        string(task.ID()),
		Username:  username,
		Title:     task.Title(),
		CreatedAt: task.CreatedAt(),
	}
}
