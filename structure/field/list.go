package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type List struct {
	*Field
}

func (text *List) ComponentGroup() string {
	return "field"
}

func NewList() component.Interface {
	return &List{
		Field: &Field{},
	}
}
