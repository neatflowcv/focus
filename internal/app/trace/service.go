package trace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	repo repository.TraceRepository
}

func NewService(repo repository.TraceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrace(ctx context.Context, input *CreateTraceInput) error {
	trace := domain.NewTrace(
		domain.TraceID(input.ID),
		domain.TraceID(input.ParentID),
		time.Duration(0),
		time.Duration(0),
		time.Time{},
	)

	err := s.repo.CreateTrace(ctx, trace)
	if err != nil {
		return fmt.Errorf("failed to create trace: %w", err)
	}

	return nil
}

func (s *Service) DeleteTrace(ctx context.Context, input *DeleteTraceInput) error {
	trace, err := s.repo.GetTrace(ctx, domain.TraceID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrTraceNotFound) {
			return ErrTraceNotFound
		}

		return fmt.Errorf("failed to get trace: %w", err)
	}

	var updates []*domain.Trace

	parentID := trace.ParentID()
	for parentID != "" {
		parent, err := s.repo.GetTrace(ctx, parentID)
		if err != nil {
			return fmt.Errorf("failed to get trace: %w", err)
		}

		update := parent.SetActual(parent.Actual() - trace.Actual())
		updates = append(updates, update)

		parentID = parent.ParentID()
	}

	err = s.repo.UpdateTraces(ctx, updates...)
	if err != nil {
		return fmt.Errorf("failed to update trace: %w", err)
	}

	err = s.repo.DeleteTrace(ctx, trace)
	if err != nil {
		return fmt.Errorf("failed to delete trace: %w", err)
	}

	return nil
}

func (s *Service) SetActual(ctx context.Context, input *SetActualInput) error {
	trace, err := s.repo.GetTrace(ctx, domain.TraceID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrTraceNotFound) {
			return ErrTraceNotFound
		}
	}

	var updates []*domain.Trace

	update := trace.SetActual(input.Actual)
	updates = append(updates, update)

	diff := input.Actual - trace.Actual()

	parentID := trace.ParentID()
	for parentID != "" {
		parent, err := s.repo.GetTrace(ctx, parentID)
		if err != nil {
			return fmt.Errorf("failed to get trace: %w", err)
		}

		update = parent.SetActual(parent.Actual() + diff)
		updates = append(updates, update)

		parentID = parent.ParentID()
	}

	err = s.repo.UpdateTraces(ctx, updates...)
	if err != nil {
		return fmt.Errorf("failed to update trace: %w", err)
	}

	return nil
}

func (s *Service) UpdateParent(ctx context.Context, input *UpdateParentInput) error {
	panic("not implemented")
}

func (s *Service) ListTraces(ctx context.Context, input *ListTracesInput) (*ListTracesOutput, error) {
	var ids []domain.TraceID
	for _, id := range input.IDs {
		ids = append(ids, domain.TraceID(id))
	}

	traces, err := s.repo.ListTraces(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to list traces: %w", err)
	}

	var items []*Trace
	for _, trace := range traces {
		items = append(items, &Trace{
			Estimated: trace.Estimated(),
			Actual:    trace.Actual(),
			StartedAt: trace.StartedAt(),
		})
	}

	return &ListTracesOutput{
		Traces: items,
	}, nil
}
