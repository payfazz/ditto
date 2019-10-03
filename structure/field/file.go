package field

import (
	"github.com/payfazz/ditto/structure/component"
)

type File struct {
	*Field
}

func (text *File) ComponentGroup() string {
	return "field"
}

func NewFile() component.Interface {
	return &File{
		Field: &Field{},
	}
}
