package memory

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	tasks map[string]map[domain.TaskID]*domain.Task
}

func NewRepository() *Repository {
	return &Repository{
		tasks: make(map[string]map[domain.TaskID]*domain.Task),
	}
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	if _, ok := r.tasks[username]; !ok {
		r.tasks[username] = make(map[domain.TaskID]*domain.Task)
	}

	r.tasks[username][task.ID()] = task

	return nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error) {
	task, ok := r.tasks[username][id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	return task, nil
}
