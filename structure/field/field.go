package field

import "github.com/payfazz/ditto/structure/component"

type Field struct {
	comp *component.Component
}

func (field *Field) RequiredAttrs() []string {
	parentKeys := field.comp.RequiredAttrs()
	return append(parentKeys, "validations")
}
