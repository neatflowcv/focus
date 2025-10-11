package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/extra"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/app/trace"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ task.Service = (*Handler)(nil)

type Handler struct {
	flowService  *flow.Service
	extraService *extra.Service
	traceService *trace.Service
	vault        *vault.Vault
}

func NewHandler(flowService *flow.Service, extraService *extra.Service, traceService *trace.Service) *Handler {
	return &Handler{
		flowService:  flowService,
		extraService: extraService,
		traceService: traceService,
		vault:        vault.NewVault("key-stone", []byte("asdf")),
	}
}

func (h *Handler) Setup(ctx context.Context, input *task.SetupTaskInput) error {
	username, _, err := h.authUser(input.Authorization)
	if err != nil {
		return err
	}

	err = h.flowService.CreateRootDummy(ctx, &flow.CreateRootDummyInput{
		Username: username,
	})
	if err != nil {
		return task.MakeInternalServerError(err)
	}

	return nil
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

	out, err := h.flowService.CreateTask(ctx, &flow.CreateTaskInput{
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

	ret, err := h.flowService.ListTasks(ctx, &flow.ListTasksInput{
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

	err = h.flowService.DeleteTask(ctx, &flow.DeleteTaskInput{
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
