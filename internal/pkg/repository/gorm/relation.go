package gorm

import (
	"database/sql"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Relation struct {
	ID       string `gorm:"primaryKey"`
	ParentID sql.NullString
	NextID   sql.NullString
	Version  uint64
}

func (r *Relation) ToDomain() *domain.Relation {
	return domain.NewRelation(
		domain.RelationID(r.ID),
		domain.RelationID(r.ParentID.String),
		domain.RelationID(r.NextID.String),
		r.Version,
	)
}

func FromDomainRelation(relation *domain.Relation) *Relation {
	return &Relation{
		ID:       string(relation.ID()),
		ParentID: sql.NullString{String: string(relation.ParentID()), Valid: true},
		NextID:   sql.NullString{String: string(relation.NextID()), Valid: true},
		Version:  relation.Version(),
	}
}

func ToDomainRelations(relations []Relation) []*domain.Relation {
	var ret []*domain.Relation
	for _, relation := range relations {
		ret = append(ret, relation.ToDomain())
	}

	return ret
}
