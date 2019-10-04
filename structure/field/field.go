package field

import (
	"errors"
	"fmt"
	"github.com/payfazz/ditto/structure/component"
)

type Field struct {
	*component.Component
	Validations []Validation
}

func NewField() *Field {
	return &Field{
		Component: &component.Component{},
	}
}

type Validation struct {
	Type         string `json:"type"`
	ErrorMessage string `json:"error_message"`
	Value        string `json:"value,omitempty"`
}

func (field *Field) RequiredAttrs() []string {
	parentKeys := field.Component.RequiredAttrs()
	return append(parentKeys, "validations")
}

func (field *Field) FillStruct(attrs component.Attributes) error {
	err := component.ValidateRequiredAttrs(field, attrs)
	if nil != err {
		return err
	}

	err = field.Component.FillStruct(attrs)
	if nil != err {
		return err
	}

	vals, ok := attrs["validations"].([]interface{})
	fmt.Println(vals)
	if !ok {
		return errors.New("validation should be array")
	}

	validations, err := extractValidations(vals)
	if nil != err {
		return err
	}

	field.Validations = validations
	return nil
}

func extractValidations(vals []interface{}) ([]Validation, error) {
	result := make([]Validation, 0)
	for _, val := range vals {
		data, ok := val.(map[string]interface{})
		if !ok {
			return nil, errors.New("validation must be object")
		}

		v := Validation{}

		if data["error_message"] == nil {
			return nil, errors.New(`validation should have property: error_message`)
		}

		v.ErrorMessage = data["error_message"].(string)

		if data["type"] == nil {
			return nil, errors.New(`validation should have property: type`)
		}

		v.Type = data["type"].(string)

		if data["value"] != nil {
			v.Value = data["value"].(string)
		}
		result = append(result, v)
	}

	return result, nil
}
