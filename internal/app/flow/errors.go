package flow

import "errors"

var (
	ErrParentTaskNotFound = errors.New("parent task not found")
	ErrTaskNotFound       = errors.New("task not found")
)
