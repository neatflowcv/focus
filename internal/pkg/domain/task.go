package domain

import "time"

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

type Task struct {
	id        string
	parentID  string
	title     string
	createdAt time.Time
	status    TaskStatus
	order     float64
}

func NewTask(id, parentID, title string, createdAt time.Time, status TaskStatus, order float64) *Task {
	return &Task{
		id:        id,
		parentID:  parentID,
		title:     title,
		createdAt: createdAt,
		status:    status,
		order:     order,
	}
}

func (t *Task) ID() string {
	return t.id
}

func (t *Task) ParentID() string {
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
