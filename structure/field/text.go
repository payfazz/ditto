package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type Text struct {
	*Field
}

func (text *Text) ComponentGroup() string {
	return "field"
}

func NewText() component.Interface {
	return &Text{
		Field: &Field{},
	}
}

