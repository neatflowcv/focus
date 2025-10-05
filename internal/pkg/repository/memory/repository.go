package memory

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	tasks    map[string]map[string]*domain.Task
	children map[string]map[string][]*domain.Task
}

func NewRepository() *Repository {
	return &Repository{
		tasks:    make(map[string]map[string]*domain.Task),
		children: make(map[string]map[string][]*domain.Task),
	}
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	if _, ok := r.tasks[username]; !ok {
		r.tasks[username] = make(map[string]*domain.Task)
	}

	if _, ok := r.children[username]; !ok {
		r.children[username] = make(map[string][]*domain.Task)
	}

	r.tasks[username][task.ID()] = task
	r.children[username][task.ParentID()] = append(r.children[username][task.ParentID()], task)

	return nil
}

func (r *Repository) ListTasks(
	ctx context.Context,
	username string,
	parentID string,
	recursive bool,
) ([]*domain.Task, error) {
	if recursive {
		var (
			tasks []*domain.Task
			stack []string
		)

		stack = append(stack, parentID)
		for len(stack) > 0 {
			parentID = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, task := range r.children[username][parentID] {
				tasks = append(tasks, task)
				stack = append(stack, task.ID())
			}
		}

		return tasks, nil
	}

	return r.children[username][parentID], nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id string) (*domain.Task, error) {
	task, ok := r.tasks[username][id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	return task, nil
}

func (r *Repository) CountSubtasks(ctx context.Context, username string, id string) (int, error) {
	return len(r.children[username][id]), nil
}
