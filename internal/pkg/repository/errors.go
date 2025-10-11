package repository

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskBusy          = errors.New("task busy")
)

var (
	ErrExtraNotFound      = errors.New("extra not found")
	ErrExtraAlreadyExists = errors.New("extra already exists")
)

var (
	ErrTraceNotFound      = errors.New("trace not found")
	ErrTraceAlreadyExists = errors.New("trace already exists")
)
