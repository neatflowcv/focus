package flow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/idmaker"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	bus     *eventbus.Bus
	idmaker idmaker.IDMaker
	repo    repository.Repository
}

func NewService(bus *eventbus.Bus, idmaker idmaker.IDMaker, repo repository.Repository) *Service {
	return &Service{
		bus:     bus,
		idmaker: idmaker,
		repo:    repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, input *CreateTaskInput) (*domain.Task, error) {
	task := domain.NewTask(
		domain.TaskID(s.idmaker.MakeID()),
		input.Title,
		domain.TaskStatusTodo,
		input.Now,
		time.Time{},
	)

	err := s.repo.CreateTask(ctx, input.Username, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	s.bus.TaskCreated.Publish(ctx, &eventbus.TaskCreatedEvent{
		TaskID: string(task.ID()),
	})

	return task, nil
}

func (s *Service) ListTasks(ctx context.Context, input *ListTasksInput) ([]*domain.Task, error) {
	var ret []*domain.Task

	for _, id := range input.IDs {
		task, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(id))
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				continue
			}

			return nil, fmt.Errorf("failed to get task: %w", err)
		}

		ret = append(ret, task)
	}

	return ret, nil
}
