package repository

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Repository interface {
	CreateTask(ctx context.Context, username string, task *domain.Task) error
	GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error)
}
