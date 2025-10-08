package main

import (
	"context"
	"strings"
	"time"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ task.Service = (*Handler)(nil)

type Handler struct {
	service *flow.Service
	vault   *vault.Vault
}

func NewHandler(service *flow.Service) *Handler {
	return &Handler{service: service, vault: vault.NewVault("key-stone", []byte("asdf"))}
}

func (h *Handler) Create(ctx context.Context, input *task.TaskInput) (*task.Taskdetail, error) {
	token := strings.TrimPrefix(input.Authorization, "Bearer ")
	now := time.Now()

	username, err := h.vault.Decrypt(token, now)
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	ret, err := h.service.CreateTask(ctx, &flow.CreateTaskInput{
		Username: username,
		Title:    input.Title,
		Now:      now,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return toTaskDetail(ret), nil
}

func (h *Handler) List(ctx context.Context, input *task.ListPayload) (task.TaskdetailCollection, error) {
	panic("unimplemented")
}

func toTaskDetail(domainTask *domain.Task) *task.Taskdetail {
	return &task.Taskdetail{
		ID:            string(domainTask.ID()),
		ParentID:      nil,
		Title:         domainTask.Title(),
		CreatedAt:     domainTask.CreatedAt().Unix(),
		Status:        string(domainTask.Status()),
		Order:         0.0,
		IsLeaf:        nil,
		CompletedAt:   nil,
		StartedAt:     nil,
		LeadTime:      nil,
		EstimatedTime: nil,
		ActualTime:    nil,
	}
}

func (h *Handler) Delete(context.Context, *task.TaskDeleteInput) error {
	panic("unimplemented")
}

func (h *Handler) Update(context.Context, *task.TaskUpdateInput) (*task.Taskdetail, error) {
	panic("unimplemented")
}
