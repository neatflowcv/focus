package extra

import "errors"

var (
	ErrPreconditionFailed = errors.New("precondition failed")
	ErrExtraNotFound      = errors.New("extra not found")
)
