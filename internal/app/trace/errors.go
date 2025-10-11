package trace

import "errors"

var (
	ErrTraceNotFound       = errors.New("trace not found")
	ErrParentTraceNotFound = errors.New("parent trace not found")
)
