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

func (s *Service) CreateTask(ctx context.Context, input *CreateTaskInput) (*CreateTaskOutput, error) {
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

	return &CreateTaskOutput{
		Task: Task{
			ID:        string(task.ID()),
			Title:     task.Title(),
			CreatedAt: task.CreatedAt(),
			Status:    string(task.Status()),
		},
	}, nil
}

func (s *Service) ListTasks(ctx context.Context, input *ListTasksInput) (*ListTasksOutput, error) {
	var tasks []Task

	for _, id := range input.IDs {
		task, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(id))
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				continue
			}

			return nil, fmt.Errorf("failed to get task: %w", err)
		}

		tasks = append(tasks, Task{
			ID:        string(task.ID()),
			Title:     task.Title(),
			CreatedAt: task.CreatedAt(),
			Status:    string(task.Status()),
		})
	}

	return &ListTasksOutput{
		Tasks: tasks,
	}, nil
}

func (s *Service) DeleteTask(ctx context.Context, input *DeleteTaskInput) error {
	task, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(input.TaskID))
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to get task: %w", err)
	}

	err = s.repo.DeleteTask(ctx, input.Username, task)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	s.bus.TaskDeleted.Publish(ctx, &eventbus.TaskDeletedEvent{
		TaskID: string(task.ID()),
	})

	return nil
}
