package repository

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type RelationRepository interface {
	CreateRelation(ctx context.Context, relation *domain.Relation) error
	DeleteRelation(ctx context.Context, relation *domain.Relation) error
	UpdateRelations(ctx context.Context, relations ...*domain.Relation) error
	GetRelation(ctx context.Context, id domain.RelationID) (*domain.Relation, error)
	ListChildrenRelations(ctx context.Context, id domain.RelationID) ([]*domain.Relation, error)
}
