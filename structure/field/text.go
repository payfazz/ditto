package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type Text struct {
	*Field
}

func NewText() component.Interface {
	return &Text{
		Field: &Field{},
	}
}

