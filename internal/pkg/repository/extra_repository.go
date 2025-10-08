package repository

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type ExtraRepository interface {
	CreateExtra(ctx context.Context, extra *domain.Extra) error
	DeleteExtra(ctx context.Context, extra *domain.Extra) error
	UpdateExtra(ctx context.Context, extra *domain.Extra) error
	GetExtra(ctx context.Context, id domain.ExtraID) (*domain.Extra, error)
	ListExtras(ctx context.Context, ids []domain.ExtraID) ([]*domain.Extra, error)
}
