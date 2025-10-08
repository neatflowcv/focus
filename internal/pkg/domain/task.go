package domain

import "time"

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

type TaskID string

type Task struct {
	id     TaskID
	title  string
	status TaskStatus

	createdAt   time.Time
	completedAt time.Time
}

func NewTask(
	id TaskID,
	title string,
	status TaskStatus,
	createdAt time.Time,
	completedAt time.Time,
) *Task {
	return &Task{
		id:          id,
		title:       title,
		status:      status,
		createdAt:   createdAt,
		completedAt: completedAt,
	}
}

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) CompletedAt() time.Time {
	return t.completedAt
}
