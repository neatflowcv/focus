package ulid

import (
	"github.com/neatflowcv/focus/internal/pkg/idmaker"
	"github.com/oklog/ulid/v2"
)

var _ idmaker.IDMaker = &IDMaker{}

type IDMaker struct{}

func NewIDMaker() *IDMaker {
	return &IDMaker{}
}

func (i *IDMaker) MakeID() string {
	return ulid.Make().String()
}
