package domain

import (
	"time"
)

type TaskID string

type Task struct {
	id        TaskID
	title     string
	createdAt time.Time
}

func NewTask(
	id TaskID,
	title string,
	createdAt time.Time,
) *Task {
	if id == "" {
		panic("id is required")
	}

	if title == "" {
		panic("title is required")
	}

	return &Task{
		id:        id,
		title:     title,
		createdAt: createdAt,
	}
}

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}
