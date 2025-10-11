package flow

import (
	"context"
	"errors"
	"fmt"

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
	if input.ParentID != "" {
		_, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(input.ParentID))
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				return nil, ErrParentTaskNotFound
			}

			return nil, fmt.Errorf("failed to get parent task: %w", err)
		}
	}

	if input.NextID != "" {
		_, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(input.NextID))
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				return nil, ErrNextTaskNotFound
			}

			return nil, fmt.Errorf("failed to get next task: %w", err)
		}
	}

	task := domain.NewTask(
		domain.TaskID(s.idmaker.MakeID()),
		domain.TaskID(input.ParentID),
		domain.TaskID(input.NextID),
		input.Title,
		input.Now,
		1,
	)
	dummy := task.Dummy()

	err := s.repo.CreateTasks(ctx, input.Username, task, dummy)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	s.bus.TaskCreated.Publish(ctx, &eventbus.TaskCreatedEvent{
		TaskID: string(task.ID()),
	})

	return &CreateTaskOutput{
		ID:        string(task.ID()),
		CreatedAt: task.CreatedAt(),
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

func (s *Service) ListTasks(ctx context.Context, input *ListTasksInput) (*ListTasksOutput, error) {
	tasks, err := s.repo.ListTasks(ctx, input.Username, domain.TaskID(input.ParentID))
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	var items []*Task
	for _, task := range tasks {
		items = append(items, &Task{
			ID:        string(task.ID()),
			Title:     task.Title(),
			CreatedAt: task.CreatedAt(),
		})
	}

	return &ListTasksOutput{
		Tasks: items,
	}, nil
}

func (s *Service) GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
	task, err := s.repo.GetTask(ctx, input.Username, domain.TaskID(input.TaskID))
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &GetTaskOutput{
		Task: Task{
			ID:        string(task.ID()),
			Title:     task.Title(),
			CreatedAt: task.CreatedAt(),
		},
	}, nil
}
