package gorm

import "github.com/neatflowcv/focus/internal/pkg/domain"

type TaskRelation struct {
	TaskID   string `gorm:"primaryKey"`
	ParentID string
}

func FromDomainTaskRelation(task *domain.Task) *TaskRelation {
	return &TaskRelation{
		TaskID:   string(task.ID()),
		ParentID: string(task.ParentID()),
	}
}
