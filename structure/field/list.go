package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type List struct {
	*Field
}

func NewList() component.Interface {
	return &List{
		Field: NewField(),
	}
}
