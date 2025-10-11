package flow

import (
	"time"
)

type CreateTaskInput struct {
	Username string
	Title    string
	ParentID string
	NextID   string
	Now      time.Time
}

type CreateTaskOutput struct {
	ID        string
	CreatedAt time.Time
	Version   uint64
}

type ListTasksInput struct {
	Username string
	ParentID string
}

type Task struct {
	ID        string
	Title     string
	CreatedAt time.Time
}

type ListTasksOutput struct {
	Tasks []*Task
}

type DeleteTaskInput struct {
	Username string
	TaskID   string
}

type GetTaskInput struct {
	Username string
	TaskID   string
}

type GetTaskOutput struct {
	Task Task
}

type CreateRootDummyInput struct {
	Username string
}

type UpdateTaskInput struct {
	Username string
	TaskID   string
	ParentID string // Next을 가져와 ParentID를 가져올 수는 있으나, NextID가 "" 인 경우를 위해 필요함
	NextID   string
	Title    string
}
