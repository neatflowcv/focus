package gorm

import (
	"database/sql"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Task struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	ParentID  sql.NullString
	Title     string
	CreatedAt time.Time
	Status    string
	Order     float64
}

func (t *Task) ToDomain() *domain.Task {
	return domain.NewTask(
		t.ID,
		t.ParentID.String,
		t.Title,
		t.CreatedAt,
		domain.TaskStatus(t.Status),
		t.Order,
	)
}

func FromDomainTask(domainTask *domain.Task, username string) *Task {
	return &Task{
		ID:        domainTask.ID(),
		ParentID:  sql.NullString{String: domainTask.ParentID(), Valid: domainTask.ParentID() != ""},
		Title:     domainTask.Title(),
		CreatedAt: domainTask.CreatedAt(),
		Status:    string(domainTask.Status()),
		Order:     domainTask.Order(),
		Username:  username,
	}
}
