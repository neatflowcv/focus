package main

import (
	"context"
	"errors"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/flow"
)

var _ task.Service = (*Handler)(nil)

type Handler struct {
	service *flow.Service
}

func NewHandler(service *flow.Service) *Handler {
	return &Handler{service: service}
}

var errNotImplemented = errors.New("not implemented")

func (h *Handler) Create(context.Context, *task.TaskInput) (*task.Taskdetail, error) {
	return nil, errNotImplemented
}

func (h *Handler) List(context.Context, *task.ListPayload) (task.TaskdetailCollection, error) {
	return nil, errNotImplemented
}
