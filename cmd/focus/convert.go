package main

import (
	"time"

	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/extra"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/app/trace"
)

func makeCreateTaskOutput(
	in *task.CreateTaskInput,
	out *flow.CreateTaskOutput,
	extraOut *extra.ListExtrasOutput,
	traceOut *trace.ListTracesOutput,
) *task.Createtaskoutput {
	return &task.Createtaskoutput{
		ID:            out.ID,
		ParentID:      in.ParentID,
		Title:         in.Title,
		CreatedAt:     out.CreatedAt.Unix(),
		EstimatedTime: pointer(int64(traceOut.Traces[0].Estimated.Seconds())),
		ActualTime:    pointer(int64(traceOut.Traces[0].Actual.Seconds())),
		StartedAt:     startedAt(traceOut.Traces[0].StartedAt),
		Status:        &extraOut.Extras[0].Status,
		IsLeaf:        &extraOut.Extras[0].Leaf,
	}
}

func pointer[T any](v T) *T {
	return &v
}

func makeCreatetaskoutputCollection(
	in *task.ListPayload,
	flowOut *flow.ListTasksOutput,
	extraOut *extra.ListExtrasOutput,
	traceOut *trace.ListTracesOutput,
) task.CreatetaskoutputCollection {
	var ret task.CreatetaskoutputCollection
	for idx, item := range flowOut.Tasks {
		ret = append(ret, &task.Createtaskoutput{
			ID:            item.ID,
			ParentID:      in.ParentID,
			Title:         item.Title,
			CreatedAt:     item.CreatedAt.Unix(),
			EstimatedTime: pointer(int64(traceOut.Traces[idx].Estimated.Seconds())),
			ActualTime:    pointer(int64(traceOut.Traces[idx].Actual.Seconds())),
			StartedAt:     startedAt(traceOut.Traces[idx].StartedAt),
			Status:        &extraOut.Extras[idx].Status,
			IsLeaf:        &extraOut.Extras[idx].Leaf,
		})
	}

	return ret
}

func makeUpdateTaskOutput(
	in *task.TaskUpdateInput,
	flowOut *flow.GetTaskOutput,
	extraOut *extra.ListExtrasOutput,
	traceOut *trace.ListTracesOutput,
) *task.Createtaskoutput {
	return &task.Createtaskoutput{
		ID:            in.TaskID,
		ParentID:      in.ParentID,
		Title:         in.Title,
		CreatedAt:     flowOut.Task.CreatedAt.Unix(),
		EstimatedTime: pointer(int64(traceOut.Traces[0].Estimated.Seconds())),
		ActualTime:    pointer(int64(traceOut.Traces[0].Actual.Seconds())),
		StartedAt:     startedAt(traceOut.Traces[0].StartedAt),
		Status:        &extraOut.Extras[0].Status,
		IsLeaf:        &extraOut.Extras[0].Leaf,
	}
}

func startedAt(t time.Time) *int64 {
	if t.IsZero() {
		return nil
	}

	return pointer(t.Unix())
}
