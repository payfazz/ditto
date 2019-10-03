package field

import "github.com/payfazz/ditto/structure/component"

type Field struct {
	comp *component.Component
}

func (field *Field) RequiredKeys() []string {
	parentKeys := field.comp.RequiredKeys()
	return append(parentKeys, "validations")
}
