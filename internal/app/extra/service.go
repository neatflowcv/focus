package extra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/focus/internal/app/relation"
	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	repo            repository.ExtraRepository
	relationService *relation.Service
}

func NewService(repo repository.ExtraRepository, relationService *relation.Service) *Service {
	return &Service{
		repo:            repo,
		relationService: relationService,
	}
}

func (s *Service) CreateExtra(ctx context.Context, input *CreateExtraInput) (*CreateExtraOutput, error) {
	extra := domain.NewExtra(
		domain.ExtraID(input.ID),
		time.Duration(0),
		time.Duration(0),
		time.Time{},
		true,
	)

	err := s.repo.CreateExtra(ctx, extra)
	if err != nil {
		return nil, fmt.Errorf("failed to create extra: %w", err)
	}

	searchID := input.ID
	for searchID != "" {
		rel, err := s.relationService.GetRelation(ctx, &relation.GetRelationInput{
			ID: searchID,
		})
		if err != nil {
			if errors.Is(err, relation.ErrRelationNotFound) {
				return nil, ErrPreconditionFailed
			}

			return nil, fmt.Errorf("failed to get relation: %w", err)
		}

		parentExtra, err := s.repo.GetExtra(ctx, domain.ExtraID(rel.ParentID))
		if err != nil {
			return nil, fmt.Errorf("failed to get parent extra: %w", err)
		}

		if !parentExtra.Leaf() {
			// parent가 이미 leaf가 아니라면, 위에 있는 모든 extra도 leaf가 아니므로 종료
			break
		}

		update := parentExtra.SetLeaf(false)

		err = s.repo.UpdateExtra(ctx, update)
		if err != nil {
			return nil, fmt.Errorf("failed to update parent extra: %w", err)
		}

		searchID = rel.ParentID
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
	return &ListExtrasOutput{
		Extras: []Extra{},
	}, nil
}

func (s *Service) UpdateActualTime(ctx context.Context, input *UpdateActualTimeInput) error {
	searchID := input.ID

	for searchID != "" {
		extra, err := s.repo.GetExtra(ctx, domain.ExtraID(searchID))
		if err != nil {
			return fmt.Errorf("failed to get extra: %w", err)
		}

		update := extra.SetActualTime(extra.ActualTime() + input.ActualTime)

		err = s.repo.UpdateExtra(ctx, update)
		if err != nil {
			return fmt.Errorf("failed to update extra: %w", err)
		}

		rel, err := s.relationService.GetRelation(ctx, &relation.GetRelationInput{
			ID: searchID,
		})
		if err != nil {
			if errors.Is(err, relation.ErrRelationNotFound) {
				return ErrPreconditionFailed
			}

			return fmt.Errorf("failed to get relation: %w", err)
		}

		searchID = rel.ParentID
	}

	return nil
}
