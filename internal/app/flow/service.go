package flow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/idmaker"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	idmaker idmaker.IDMaker
	repo    repository.Repository
}

func NewService(idmaker idmaker.IDMaker, repo repository.Repository) *Service {
	return &Service{idmaker: idmaker, repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, input *CreateTaskInput) (*domain.Task, error) {
	if input.ParentID != "" {
		_, err := s.repo.GetTask(ctx, input.Username, input.ParentID)
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				return nil, ErrParentTaskNotFound
			}

			return nil, fmt.Errorf("failed to get parent task: %w", err)
		}
	}

	count, err := s.repo.CountSubtasks(ctx, input.Username, input.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to count subtasks: %w", err)
	}

	task := domain.NewTask(
		domain.TaskID(s.idmaker.MakeID()),
		input.ParentID,
		input.Title,
		input.Now,
		domain.TaskStatusTodo,
		float64(count)*10.0+10.0, //nolint:mnd
		time.Time{},
		time.Time{},
		time.Duration(0),
		time.Duration(0),
	)

	err = s.repo.CreateTask(ctx, input.Username, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (s *Service) ListTasks(ctx context.Context, input *ListTasksInput) ([]*domain.Task, error) {
	var listFn func(ctx context.Context, username string, parentID domain.TaskID) ([]*domain.Task, error)
	if input.Recursive {
		listFn = s.repo.ListDescendantsTasks
	} else {
		listFn = s.repo.ListSubTasks
	}

	ret, err := listFn(ctx, input.Username, input.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return ret, nil
}
