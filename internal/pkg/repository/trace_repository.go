package repository

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type TraceRepository interface {
	CreateTrace(ctx context.Context, trace *domain.Trace) error
	DeleteTrace(ctx context.Context, trace *domain.Trace) error
	GetTrace(ctx context.Context, id domain.TraceID) (*domain.Trace, error)
	UpdateTraces(ctx context.Context, traces ...*domain.Trace) error
	ListTraces(ctx context.Context, ids []domain.TraceID) ([]*domain.Trace, error)
	ListChildTraces(ctx context.Context, parentID domain.TraceID) ([]*domain.Trace, error)
}
