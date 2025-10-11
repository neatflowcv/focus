package trace

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

func (s *Service) UpdateParent(ctx context.Context, input *UpdateParentInput) error { //nolint:cyclop,funlen
	trace, err := s.repo.GetTrace(ctx, domain.TraceID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrTraceNotFound) {
			return ErrTraceNotFound
		}

		return fmt.Errorf("failed to get trace: %w", err)
	}

	if input.ParentID != "" {
		_, err := s.repo.GetTrace(ctx, domain.TraceID(input.ParentID))
		if err != nil {
			if errors.Is(err, repository.ErrTraceNotFound) {
				return ErrParentTraceNotFound
			}
		}
	}

	score := s.taskScore(ctx, trace)

	var updates []*domain.Trace

	update := trace.SetParentID(domain.TraceID(input.ParentID))
	updates = append(updates, update)

	if score > 0 {
		oldParents, err := s.findAncestors(ctx, trace.ParentID())
		if err != nil {
			return err
		}

		newParents, err := s.findAncestors(ctx, update.ParentID())
		if err != nil {
			return err
		}

		pivot := 0
		for pivot < len(oldParents) && pivot < len(newParents) {
			if oldParents[pivot].ID() != newParents[pivot].ID() {
				break
			}

			pivot++
		}

		idx := pivot

		for idx < len(oldParents) {
			oldUpdate := oldParents[idx].SetActual(oldParents[idx].Actual() - score)
			updates = append(updates, oldUpdate)
			idx++
		}

		idx = pivot
		for idx < len(newParents) {
			newUpdate := newParents[idx].SetActual(newParents[idx].Actual() + score)
			updates = append(updates, newUpdate)
			idx++
		}
	}

	err = s.repo.UpdateTraces(ctx, updates...)
	if err != nil {
		return fmt.Errorf("failed to update trace: %w", err)
	}

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, input *UpdateStatusInput) error {
	trace, err := s.repo.GetTrace(ctx, domain.TraceID(input.ID))
	if err != nil {
		return fmt.Errorf("failed to get trace: %w", err)
	}

	switch domain.TaskStatus(input.Status) {
	case domain.TaskStatusDoing:
		err := s.startTrace(ctx, trace, input.Now)
		if err != nil {
			return err
		}

	case domain.TaskStatusDone, domain.TaskStatusTodo:
		err := s.stopTrace(ctx, trace, input.Now)
		if err != nil {
			return err
		}
	}

	return nil
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

func (s *Service) startTrace(ctx context.Context, trace *domain.Trace, now time.Time) error {
	if !trace.StartedAt().IsZero() {
		// already started
		return nil
	}

	update := trace.SetStartedAt(now)

	err := s.repo.UpdateTraces(ctx, update)
	if err != nil {
		return fmt.Errorf("failed to update trace: %w", err)
	}

	return nil
}

func (s *Service) stopTrace(ctx context.Context, trace *domain.Trace, now time.Time) error {
	if trace.StartedAt().IsZero() {
		// already stopped
		return nil
	}

	diff := now.Sub(trace.StartedAt())

	ancestors, err := s.findAncestors(ctx, trace.ParentID())
	if err != nil {
		return err
	}

	ancestors = append(ancestors, trace)

	var updates []*domain.Trace

	for _, ancestor := range ancestors {
		update := ancestor.
			SetStartedAt(time.Time{}).
			SetActual(ancestor.Actual() + diff)
		updates = append(updates, update)
	}

	err = s.repo.UpdateTraces(ctx, updates...)
	if err != nil {
		return fmt.Errorf("failed to update trace: %w", err)
	}

	return nil
}

func (s *Service) taskScore(ctx context.Context, trace *domain.Trace) time.Duration {
	traces, err := s.repo.ListChildTraces(ctx, trace.ID())
	if err != nil {
		return 0
	}

	sum := time.Duration(0)
	for _, trace := range traces {
		sum += trace.Actual()
	}

	return trace.Actual() - sum
}

func (s *Service) findAncestors(ctx context.Context, id domain.TraceID) ([]*domain.Trace, error) {
	var parents []*domain.Trace

	searchID := id
	for searchID != "" {
		parent, err := s.repo.GetTrace(ctx, searchID)
		if err != nil {
			if errors.Is(err, repository.ErrTraceNotFound) {
				return nil, ErrParentTraceNotFound
			}

			return nil, fmt.Errorf("failed to get trace: %w", err)
		}

		parents = append(parents, parent)

		searchID = parent.ParentID()
	}

	slices.Reverse(parents)

	return parents, nil
}
