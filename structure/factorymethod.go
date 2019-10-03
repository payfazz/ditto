package structure

import (
	"errors"
	"github.com/payfazz/ditto/structure/component"
	"github.com/payfazz/ditto/structure/field"
)

type FactoryMethod func() component.Interface

var registeredComponents = make(map[string]FactoryMethod)

func RegisterComponent(name string, component FactoryMethod) {
	registeredComponents[name] = component
}

func CreateComponent(name string) (component.Interface, error) {
	val, ok := registeredComponents[name]

	if !ok {
		return nil, errors.New("component not found")
	}

	return val(), nil
}

func init() {
	RegisterComponent("file", field.NewFile)
	RegisterComponent("text", field.NewText)
	RegisterComponent("list", field.NewList)
}
