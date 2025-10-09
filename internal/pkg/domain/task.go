package domain

import (
	"time"
)

type TaskID string

type Task struct {
	id          TaskID
	title       string
	createdAt   time.Time
	completedAt time.Time
}

func NewTask(
	id TaskID,
	title string,
	createdAt time.Time,
	completedAt time.Time,
) *Task {
	if id == "" {
		panic("id is required")
	}

	if title == "" {
		panic("title is required")
	}

	return &Task{
		id:          id,
		title:       title,
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

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) CompletedAt() time.Time {
	return t.completedAt
}
