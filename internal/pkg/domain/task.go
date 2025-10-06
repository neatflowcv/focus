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
	id       TaskID
	parentID TaskID
	title    string
	status   TaskStatus
	order    float64

	createdAt   time.Time
	completedAt time.Time
}

func NewTask(
	id TaskID,
	parentID TaskID,
	title string,
	status TaskStatus,
	order float64,
	createdAt time.Time,
	completedAt time.Time,
) *Task {
	return &Task{
		id:          id,
		parentID:    parentID,
		title:       title,
		createdAt:   createdAt,
		status:      status,
		order:       order,
		completedAt: completedAt,
	}
}

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) ParentID() TaskID {
	return t.parentID
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func (t *Task) Order() float64 {
	return t.order
}

func (t *Task) CompletedAt() time.Time {
	return t.completedAt
}
