package flow

import (
	"time"
)

type CreateTaskInput struct {
	Username string
	Title    string
	Now      time.Time
}

type Task struct {
	ID        string
	Title     string
	CreatedAt time.Time
	Status    string
}

type CreateTaskOutput struct {
	Task Task
}

type ListTasksInput struct {
	Username string
	IDs      []string
}

type ListTasksOutput struct {
	Tasks []Task
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
