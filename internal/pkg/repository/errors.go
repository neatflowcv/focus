package repository

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
)

var (
	ErrRelationNotFound      = errors.New("relation not found")
	ErrRelationBusy          = errors.New("relation busy")
	ErrRelationAlreadyExists = errors.New("relation already exists")
)

var (
	ErrExtraNotFound      = errors.New("extra not found")
	ErrExtraAlreadyExists = errors.New("extra already exists")
)

var (
	ErrTraceNotFound      = errors.New("trace not found")
	ErrTraceAlreadyExists = errors.New("trace already exists")
)
