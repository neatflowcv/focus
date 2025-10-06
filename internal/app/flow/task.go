package flow

import (
	"time"

	"github.com/neatflowcv/focus/internal/pkg/domain"
)

type CreateTaskInput struct {
	Username string
	ParentID domain.TaskID
	Title    string
	Now      time.Time
}

type ListTasksInput struct {
	Username  string
	ParentID  domain.TaskID
	Recursive bool
}
