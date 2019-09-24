package ditto

import "errors"

type Field struct {
	ID          string
	Type        Type
	Title       string
	Description string
	Validations []FieldValidation
	Info        map[string]interface{}
}

type FieldValidation struct {
	Type         string
	ErrorMessage string
	Value        string
}

func NewFieldFromMap(data map[string]interface{}) (*Field, error) {
	if data["id"] == nil {
		return nil, errors.New(`form_should_have_property_id`)
	}
	if data["type"] == nil {
		return nil, errors.New(`form_should_have_property_type`)
	}
	if data["title"] == nil {
		return nil, errors.New(`form_should_have_property_name`)
	}
	if data["description"] == nil {
		return nil, errors.New(`form_should_have_property_description`)
	}
	if data["validations"] == nil {
		return nil, errors.New(`form_should_have_property_validation`)
	}

	validationInterface, ok := data["validations"].([]interface{})
	if !ok {
		return nil, errors.New("form_validation_should_be_array")
	}

	fieldType, ok := data["type"].(string)
	if !ok {
		return nil, errors.New(`form_type_should_be_string`)
	}

	typ := GetType(fieldType)
	if nil == typ {
		return nil, errors.New(`form_type_not_supported`)
	}

	var info map[string]interface{}
	if data["info"] != nil {
		info, ok = data["info"].(map[string]interface{})
		if !ok {
			return nil, errors.New("field_info_should_be_an_object")
		}

		err := validateInfo(info, typ)
		if err != nil {
			return nil, err
		}
	}

	validationRules := extractArrayMap(validationInterface)
	for _, rule := range validationRules {
		err := validateValidation(rule)
		if err != nil {
			return nil, err
		}
	}

	if typ.Type != "field" {
		return nil, errors.New(`field_type_not_supported`)
	}

	return &Field{
		ID:          data["id"].(string),
		Type:        *typ,
		Title:       data["title"].(string),
		Description: data["description"].(string),
		Validations: nil, //TODO: extract validations
		Info:        info,
	}, nil
}

func validateInfo(info map[string]interface{}, typ *Type) error {
	return nil
}

func validateValidation(rule map[string]interface{}) error {
	return nil
}
