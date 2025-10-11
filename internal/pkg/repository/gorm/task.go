package gorm

import (
	"database/sql"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Task struct {
	Username string

	ID       string `gorm:"primaryKey"`
	ParentID sql.NullString
	NextID   sql.NullString

	Title     string
	CreatedAt time.Time

	Version uint64
}

func FromDomainTask(task *domain.Task, username string) *Task {
	return &Task{
		ID:        string(task.ID()),
		ParentID:  sql.NullString{String: string(task.ParentID()), Valid: true},
		NextID:    sql.NullString{String: string(task.NextID()), Valid: true},
		Username:  username,
		Title:     task.Title(),
		CreatedAt: task.CreatedAt(),
		Version:   task.Version(),
	}
}

func (t *Task) ToDomain() *domain.Task {
	return domain.NewTask(
		domain.TaskID(t.ID),
		domain.TaskID(getString(t.ParentID)),
		domain.TaskID(getString(t.NextID)),
		t.Title,
		t.CreatedAt,
		t.Version,
	)
}

func ToDomainTasks(tasks []Task) []*domain.Task {
	var ret []*domain.Task
	for _, task := range tasks {
		ret = append(ret, task.ToDomain())
	}

	return ret
}

func getString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}

	return s.String
}
