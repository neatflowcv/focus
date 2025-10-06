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
	id        TaskID
	parentID  TaskID
	title     string
	createdAt time.Time
	status    TaskStatus
	order     float64
}

func NewTask(id TaskID, parentID TaskID, title string, createdAt time.Time, status TaskStatus, order float64) *Task {
	return &Task{
		id:        id,
		parentID:  parentID,
		title:     title,
		createdAt: createdAt,
		status:    status,
		order:     order,
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
