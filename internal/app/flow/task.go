package flow

import (
	"time"
)

type CreateTaskInput struct {
	Username string
	Title    string
	Now      time.Time
}

type ListTasksInput struct {
	Username string
	IDs      []string
}
