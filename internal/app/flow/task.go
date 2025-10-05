package flow

import "time"

type CreateTaskInput struct {
	Username string
	ParentID string
	Title    string
	Now      time.Time
}

type ListTasksInput struct {
	Username  string
	ParentID  string
	Recursive bool
}
