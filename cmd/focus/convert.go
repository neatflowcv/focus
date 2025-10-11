package main

import (
	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/flow"
)

func makeCreateTaskOutput(out *flow.CreateTaskOutput) *task.Createtaskoutput {
	return &task.Createtaskoutput{
		ID:        out.ID,
		CreatedAt: out.CreatedAt.Unix(),
	}
}

func makeTaskDetail(domainTask *flow.Task, parentID string) *task.Taskdetail {
	return &task.Taskdetail{
		ID:            domainTask.ID,
		ParentID:      &parentID,
		Title:         domainTask.Title,
		CreatedAt:     domainTask.CreatedAt.Unix(),
		Status:        "",
		IsLeaf:        nil,
		CompletedAt:   nil,
		StartedAt:     nil,
		LeadTime:      nil,
		EstimatedTime: nil,
		ActualTime:    nil,
	}
}

func makeTaskdetailCollection(domainTask []*flow.Task, parentID string) task.TaskdetailCollection {
	var ret task.TaskdetailCollection
	for _, task := range domainTask {
		ret = append(ret, makeTaskDetail(task, parentID))
	}

	return ret
}
