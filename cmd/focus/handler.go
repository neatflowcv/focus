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

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	ret, err := h.service.CreateTask(ctx, &flow.CreateTaskInput{
		Username: username,
		ParentID: domain.TaskID(parentID),
		Title:    input.Title,
		Now:      now,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return toTaskDetail(ret), nil
}

func (h *Handler) List(ctx context.Context, input *task.ListPayload) (task.TaskdetailCollection, error) {
	token := strings.TrimPrefix(input.Authorization, "Bearer ")
	now := time.Now()

	username, err := h.vault.Decrypt(token, now)
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	recursive := false
	if input.Recursive != nil {
		recursive = *input.Recursive
	}

	ret, err := h.service.ListTasks(ctx, &flow.ListTasksInput{
		Username:  username,
		ParentID:  domain.TaskID(parentID),
		Recursive: recursive,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return toTaskCollection(ret), nil
}

func toTaskCollection(tasks []*domain.Task) task.TaskdetailCollection {
	var ret task.TaskdetailCollection
	for _, task := range tasks {
		ret = append(ret, toTaskDetail(task))
	}

	return ret
}

func toTaskDetail(domainTask *domain.Task) *task.Taskdetail {
	parentID := string(domainTask.ParentID())

	return &task.Taskdetail{
		ID:        string(domainTask.ID()),
		ParentID:  &parentID,
		Title:     domainTask.Title(),
		CreatedAt: domainTask.CreatedAt().Unix(),
		Status:    string(domainTask.Status()),
		Order:     domainTask.Order(),
	}
}
