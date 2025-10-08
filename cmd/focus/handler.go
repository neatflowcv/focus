package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/app/relation"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ task.Service = (*Handler)(nil)

type Handler struct {
	service         *flow.Service
	relationService *relation.Service
	vault           *vault.Vault
}

func NewHandler(service *flow.Service, relationService *relation.Service) *Handler {
	return &Handler{
		service:         service,
		relationService: relationService,
		vault:           vault.NewVault("key-stone", []byte("asdf")),
	}
}

func (h *Handler) Create(ctx context.Context, input *task.TaskInput) (*task.Taskdetail, error) {
	username, now, err := h.authUser(input.Authorization)
	if err != nil {
		return nil, err
	}

	ret, err := h.service.CreateTask(ctx, &flow.CreateTaskInput{
		Username: username,
		Title:    input.Title,
		Now:      now,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	err = h.relationService.CreateRelation(ctx, &relation.CreateRelationInput{
		ID:       ret.Task.ID,
		ParentID: parentID,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return toTaskDetail(&ret.Task, parentID), nil
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

	children, err := h.relationService.ListChildren(ctx, &relation.ListChildrenInput{
		ParentID: parentID,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	ret, err := h.service.ListTasks(ctx, &flow.ListTasksInput{
		Username: username,
		IDs:      children.IDs,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	return toTaskdetailCollection(ret.Tasks, parentID), nil
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
	username, _, err := h.authUser(input.Authorization)
	if err != nil {
		return nil, err
	}

	parentID := ""
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	nextID := ""
	if input.NextID != nil {
		nextID = *input.NextID
	}

	err = h.relationService.UpdateRelation(ctx, &relation.UpdateRelationInput{
		ID:       input.TaskID,
		ParentID: parentID,
		NextID:   nextID,
	})
	if err != nil {
		return nil, task.MakeInternalServerError(err)
	}

	ret, err := h.service.GetTask(ctx, &flow.GetTaskInput{
		Username: username,
		TaskID:   input.TaskID,
	})
	if err != nil {
		if errors.Is(err, flow.ErrTaskNotFound) {
			return nil, task.MakeTaskNotFound(err)
		}

		return nil, task.MakeInternalServerError(err)
	}

	return toTaskDetail(&ret.Task, parentID), nil
}

func toTaskDetail(domainTask *flow.Task, parentID string) *task.Taskdetail {
	return &task.Taskdetail{
		ID:            domainTask.ID,
		ParentID:      &parentID,
		Title:         domainTask.Title,
		CreatedAt:     domainTask.CreatedAt.Unix(),
		Status:        domainTask.Status,
		IsLeaf:        nil,
		CompletedAt:   nil,
		StartedAt:     nil,
		LeadTime:      nil,
		EstimatedTime: nil,
		ActualTime:    nil,
	}
}

func toTaskdetailCollection(domainTask []flow.Task, parentID string) task.TaskdetailCollection {
	var ret task.TaskdetailCollection
	for _, task := range domainTask {
		ret = append(ret, toTaskDetail(&task, parentID))
	}

	return ret
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
