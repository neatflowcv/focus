package gorm

import (
	"database/sql"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type Extra struct {
	ID       string
	ParentID sql.NullString
	Leaf     sql.NullBool
	Status   string
}

func FromDomainExtra(extra *domain.Extra) *Extra {
	return &Extra{
		ID:       string(extra.ID()),
		ParentID: sql.NullString{String: string(extra.ParentID()), Valid: true},
		Leaf:     sql.NullBool{Bool: extra.Leaf(), Valid: true},
		Status:   string(extra.Status()),
	}
}

func (e *Extra) ToDomain() *domain.Extra {
	return domain.NewExtra(
		domain.ExtraID(e.ID),
		domain.ExtraID(getString(e.ParentID)),
		getBool(e.Leaf),
		domain.TaskStatus(e.Status),
	)
}

func getBool(b sql.NullBool) bool {
	if !b.Valid {
		return false
	}

	return b.Bool
}
