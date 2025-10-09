package extra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	repo repository.ExtraRepository
}

func NewService(repo repository.ExtraRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateExtra(ctx context.Context, input *CreateExtraInput) (*CreateExtraOutput, error) {
	extra := domain.NewExtra(
		domain.ExtraID(input.ID),
		domain.ExtraID(input.ParentID),
		domain.NewDurations(time.Duration(0), time.Duration(0)),
		time.Time{},
		true,
	)

	err := s.repo.CreateExtra(ctx, extra)
	if err != nil {
		return nil, fmt.Errorf("failed to create extra: %w", err)
	}

	searchID := domain.ExtraID(input.ParentID)
	for searchID != "" {
		searchedExtra, err := s.repo.GetExtra(ctx, searchID)
		if err != nil {
			return nil, fmt.Errorf("failed to get extra: %w", err)
		}

		if !searchedExtra.Leaf() {
			// parent가 이미 leaf가 아니라면, 위에 있는 모든 extra도 leaf가 아니므로 종료
			break
		}

		update := searchedExtra.SetLeaf(false)

		err = s.repo.UpdateExtra(ctx, update)
		if err != nil {
			return nil, fmt.Errorf("failed to update parent extra: %w", err)
		}

		searchID = searchedExtra.ParentID()
	}

	return &CreateExtraOutput{
		Extra: Extra{
			EstimatedTime: extra.EstimatedTime(),
			ActualTime:    extra.ActualTime(),
			StartedAt:     extra.StartedAt(),
			Leaf:          extra.Leaf(),
		},
	}, nil
}

func (s *Service) DeleteExtra(ctx context.Context, input *DeleteExtraInput) error {
	extra, err := s.repo.GetExtra(ctx, domain.ExtraID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrExtraNotFound) {
			return ErrExtraNotFound
		}

		return fmt.Errorf("failed to get extra: %w", err)
	}

	err = s.repo.DeleteExtra(ctx, extra)
	if err != nil {
		return fmt.Errorf("failed to delete extra: %w", err)
	}

	return nil
}

func (s *Service) UpdateEstimatedTime(ctx context.Context, input *UpdateEstimatedTimeInput) error {
	extra, err := s.repo.GetExtra(ctx, domain.ExtraID(input.ID))
	if err != nil {
		return fmt.Errorf("failed to get extra: %w", err)
	}

	update := extra.SetEstimatedTime(input.EstimatedTime)

	err = s.repo.UpdateExtra(ctx, update)
	if err != nil {
		return fmt.Errorf("failed to update extra: %w", err)
	}

	return nil
}

func (s *Service) ListExtras(ctx context.Context, input *ListExtrasInput) (*ListExtrasOutput, error) {
	var ids []domain.ExtraID
	for _, id := range input.IDs {
		ids = append(ids, domain.ExtraID(id))
	}

	extras, err := s.repo.ListExtras(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to list extras: %w", err)
	}

	var ouputExtras []Extra
	for _, extra := range extras {
		ouputExtras = append(ouputExtras, Extra{
			EstimatedTime: extra.EstimatedTime(),
			ActualTime:    extra.ActualTime(),
			StartedAt:     extra.StartedAt(),
			Leaf:          extra.Leaf(),
		})
	}

	return &ListExtrasOutput{
		Extras: ouputExtras,
	}, nil
}

func (s *Service) UpdateActualTime(ctx context.Context, input *UpdateActualTimeInput) error {
	searchID := domain.ExtraID(input.ID)

	for searchID != "" {
		extra, err := s.repo.GetExtra(ctx, searchID)
		if err != nil {
			return fmt.Errorf("failed to get extra: %w", err)
		}

		update := extra.SetActualTime(extra.ActualTime() + input.ActualTime)

		err = s.repo.UpdateExtra(ctx, update)
		if err != nil {
			return fmt.Errorf("failed to update extra: %w", err)
		}

		searchID = extra.ParentID()
	}

	return nil
}
