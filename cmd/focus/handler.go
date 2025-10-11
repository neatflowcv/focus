package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ task.Service = (*Handler)(nil)

type Handler struct {
	service *flow.Service
	vault   *vault.Vault
}

func NewHandler(service *flow.Service) *Handler {
	return &Handler{
		service: service,
		vault:   vault.NewVault("key-stone", []byte("asdf")),
	}
}

func (h *Handler) Create(ctx context.Context, input *task.CreateTaskInput) (*task.Createtaskoutput, error) {
	username, now, err := h.authUser(input.Authorization)
	if err != nil {
		return nil, err
	}

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	out, err := h.service.CreateTask(ctx, &flow.CreateTaskInput{
		Username: username,
		Title:    input.Title,
		ParentID: parentID,
		NextID:   "",
		Now:      now,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return makeCreateTaskOutput(out), nil
}

func (h *Handler) List(ctx context.Context, input *task.ListPayload) (task.TaskdetailCollection, error) {
	username, _, err := h.authUser(input.Authorization)
	if err != nil {
		return nil, err
	}

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	ret, err := h.service.ListTasks(ctx, &flow.ListTasksInput{
		Username: username,
		ParentID: parentID,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return makeTaskdetailCollection(ret.Tasks, parentID), nil
}

func (h *Handler) Delete(ctx context.Context, input *task.TaskDeleteInput) error {
	username, _, err := h.authUser(input.Authorization)
	if err != nil {
		return err
	}

	err = h.service.DeleteTask(ctx, &flow.DeleteTaskInput{
		Username: username,
		TaskID:   input.TaskID,
	})
	if err != nil {
		if errors.Is(err, flow.ErrTaskNotFound) {
			return task.MakeTaskNotFound(err)
		}

		return task.MakeInternalServerError(err)
	}

	return nil
}

func (h *Handler) Update(ctx context.Context, input *task.TaskUpdateInput) (*task.Taskdetail, error) {
	panic("not implemented")
}

func (h *Handler) authUser(authorization string) (string, time.Time, error) {
	now := time.Now()
	token := strings.TrimPrefix(authorization, "Bearer ")

	username, err := h.vault.Decrypt(token, now)
	if err != nil {
		return "", now, task.MakeUnauthorized(err)
	}

	return username, now, nil
}
