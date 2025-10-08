package relation

import (
	"context"
	"errors"
	"fmt"
	"slices"

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

	err = s.repo.UpdateRelations(ctx, nextDummy)
	if err != nil {
		return fmt.Errorf("failed to update relation: %w", err)
	}

	return nil
}

func (s *Service) ListChildren(ctx context.Context, input *ListChildrenInput) (*ListChildrenOutput, error) {
	if input.ParentID != "" {
		_, err := s.repo.GetRelation(ctx, domain.RelationID(input.ParentID))
		if err != nil {
			if errors.Is(err, repository.ErrRelationNotFound) {
				return &ListChildrenOutput{
					IDs: nil,
				}, nil
			}

			return nil, fmt.Errorf("failed to get relation: %w", err)
		}
	}

	ids, err := s.listChildren(ctx, domain.RelationID(input.ParentID))
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, id := range ids[1:] {
		ret = append(ret, string(id))
	}

	return &ListChildrenOutput{
		IDs: ret,
	}, nil
}

func (s *Service) GetRelation(ctx context.Context, input *GetRelationInput) (*GetRelationOutput, error) {
	relation, err := s.repo.GetRelation(ctx, domain.RelationID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRelationNotFound) {
			return nil, ErrRelationNotFound
		}

		return nil, fmt.Errorf("failed to get relation: %w", err)
	}

	return &GetRelationOutput{
		ID:       string(relation.ID()),
		ParentID: string(relation.ParentID()),
	}, nil
}

func (s *Service) UpdateRelation(ctx context.Context, input *UpdateRelationInput) error {
	relation, err := s.repo.GetRelation(ctx, domain.RelationID(input.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRelationNotFound) {
			return ErrRelationNotFound
		}

		return fmt.Errorf("failed to get relation: %w", err)
	}

	update := relation.
		SetNextID(domain.RelationID(input.NextID)).
		SetParentID(domain.RelationID(input.ParentID))

	oldPrevRelation, err := s.getPrevRelation(ctx, relation.ParentID(), relation.ID())
	if err != nil {
		return err
	}

	oldPrevUpdate := oldPrevRelation.SetNextID(relation.NextID())

	newPrevRelation, err := s.getPrevRelation(ctx, update.ParentID(), update.NextID())
	if err != nil {
		return err
	}

	newPrevUpdate := newPrevRelation.SetNextID(relation.ID())

	err = s.repo.UpdateRelations(ctx, oldPrevUpdate, newPrevUpdate, update)
	if err != nil {
		return fmt.Errorf("failed to update relations: %w", err)
	}

	return nil
}

func (s *Service) getPrevRelation(ctx context.Context, parentID, nextID domain.RelationID) (*domain.Relation, error) {
	ids, err := s.listChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}

	prevID, err := extractPrevID(ids, nextID)
	if err != nil {
		return nil, err
	}

	ret, err := s.repo.GetRelation(ctx, prevID)
	if err != nil {
		if errors.Is(err, repository.ErrRelationNotFound) {
			return nil, fmt.Errorf("prev relation not found: %w", ErrRelationNotFound)
		}

		return nil, fmt.Errorf("failed to get prev relation: %w", err)
	}

	return ret, nil
}

func extractPrevID(ids []domain.RelationID, nextID domain.RelationID) (domain.RelationID, error) {
	if nextID == "" {
		return ids[len(ids)-1], nil
	}

	index := slices.IndexFunc(ids, func(id domain.RelationID) bool {
		return id == nextID
	})
	if index == -1 {
		return "", ErrPreviousRelationNotFound
	}

	return ids[index-1], nil
}

func (s *Service) listChildren(ctx context.Context, parentID domain.RelationID) ([]domain.RelationID, error) {
	relations, err := s.repo.ListChildrenRelations(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list children relations: %w", err)
	}

	idMap := map[domain.RelationID]*domain.Relation{}
	for _, relation := range relations {
		idMap[relation.ID()] = relation
	}

	first := domain.DummyRelationID(parentID)

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

	var ids []domain.RelationID
	for _, relation := range sortedRelation {
		ids = append(ids, relation.ID())
	}

	return ids, nil
}
