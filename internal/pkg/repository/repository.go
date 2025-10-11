package repository

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Repository interface {
	CreateTasks(ctx context.Context, username string, tasks ...*domain.Task) error
	DeleteTask(ctx context.Context, username string, task *domain.Task) error
	GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error)
	ListTasks(ctx context.Context, username string, parentID domain.TaskID) ([]*domain.Task, error)
	UpdateTasks(ctx context.Context, username string, tasks ...*domain.Task) error
}
