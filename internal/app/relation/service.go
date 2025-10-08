package relation

import (
	"context"
	"errors"
	"fmt"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

type Service struct {
	bus  *eventbus.Bus
	repo repository.RelationRepository
}

func NewService(bus *eventbus.Bus, repo repository.RelationRepository) *Service {
	return &Service{
		bus:  bus,
		repo: repo,
	}
}

func (s *Service) CreateRelation(ctx context.Context, input *CreateRelationInput) error {
	dummy, err := s.repo.GetRelation(ctx, domain.DummyRelationID(domain.RelationID(input.ParentID)))
	if err != nil {
		if errors.Is(err, repository.ErrRelationNotFound) {
			return ErrDummyNotFound
		}

		return fmt.Errorf("failed to get dummy: %w", err)
	}

	relation := domain.NewRelation(
		domain.RelationID(input.ID),
		domain.RelationID(input.ParentID),
		dummy.NextID(),
		1,
	)

	err = s.repo.CreateRelation(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRelationBusy):
			return ErrRelationBusy
		case errors.Is(err, repository.ErrRelationAlreadyExists):
			return ErrRelationAlreadyExists
		default:
			return fmt.Errorf("failed to create relation: %w", err)
		}
	}

	nextDummy := dummy.SetNextID(relation.ID())

	err = s.repo.UpdateRelation(ctx, nextDummy)
	if err != nil {
		return fmt.Errorf("failed to update relation: %w", err)
	}

	return nil
}

func (s *Service) ListChildren(ctx context.Context, input *ListChildrenInput) (*ListChildrenOutput, error) {
	relations, err := s.repo.ListChildrenRelations(ctx, domain.RelationID(input.ParentID))
	if err != nil {
		return nil, fmt.Errorf("failed to list children relations: %w", err)
	}

	idMap := map[domain.RelationID]*domain.Relation{}
	for _, relation := range relations {
		idMap[relation.ID()] = relation
	}

	first := domain.DummyRelationID(domain.RelationID(input.ParentID))

	var sortedRelation []*domain.Relation

	cur := first
	for {
		relation, ok := idMap[cur]
		if !ok {
			break
		}

		sortedRelation = append(sortedRelation, relation)
		cur = relation.NextID()
	}

	if len(sortedRelation) != len(relations) {
		panic("logic error: list children")
	}

	var ids []string
	for _, relation := range sortedRelation {
		ids = append(ids, string(relation.ID()))
	}

	if len(ids) == 1 {
		return &ListChildrenOutput{
			IDs: nil,
		}, nil
	}

	return &ListChildrenOutput{
		IDs: ids[1:],
	}, nil
}

func (s *Service) CreateChildDummy(ctx context.Context, input *CreateChildDummyInput) error {
	relation := domain.NewRelation(
		domain.RelationID(input.ID+"-dummy"),
		domain.RelationID(input.ID),
		domain.RelationID(""),
		1,
	)

	err := s.repo.CreateRelation(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRelationBusy):
			return ErrRelationBusy
		case errors.Is(err, repository.ErrRelationAlreadyExists):
			return ErrRelationAlreadyExists
		default:
			return fmt.Errorf("failed to create relation: %w", err)
		}
	}

	return nil
}

func (s *Service) DeleteChildDummy(ctx context.Context, input *DeleteChildDummyInput) error {
	relation, err := s.repo.GetRelation(ctx, domain.DummyRelationID(domain.RelationID(input.ID)))
	if err != nil {
		if errors.Is(err, repository.ErrRelationNotFound) {
			return ErrDummyNotFound
		}

		return fmt.Errorf("failed to get child dummy: %w", err)
	}

	err = s.repo.DeleteRelation(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRelationBusy):
			return ErrRelationBusy
		case errors.Is(err, repository.ErrRelationAlreadyExists):
			return ErrRelationAlreadyExists
		default:
			return fmt.Errorf("failed to delete child dummy: %w", err)
		}
	}

	return nil
}
