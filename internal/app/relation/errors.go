package relation

import "errors"

var (
	ErrDummyNotFound         = errors.New("dummy not found")
	ErrRelationBusy          = errors.New("relation busy")
	ErrRelationAlreadyExists = errors.New("relation already exists")
	ErrRelationNotFound      = errors.New("relation not found")
)
