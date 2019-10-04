package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type File struct {
	*Field
}

func NewFile() component.Interface {
	return &File{
		Field: NewField(),
	}
}
