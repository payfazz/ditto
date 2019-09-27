package structure

import (
	"errors"
	"fmt"
)

type Field struct {
	ID          string                 `json:"id"`
	Type        Type                   `json:"type"`
	Title       string                 `json:"title"`
	Description *string                `json:"description"`
	Placeholder *string                `json:"placeholder"`
	Validations []FieldValidation      `json:"validations"`
	Info        map[string]interface{} `json:"info,omitempty"`
	Status        map[string]interface{} `json:"status,omitempty"`
}

type FieldValidation struct {
	Type         string `json:"type"`
	ErrorMessage string `json:"error_message"`
	Value        string `json:"value,omitempty"`
}

func NewFieldFromMap(data map[string]interface{}) (*Field, error) {
	if data["id"] == nil {
		return nil, errors.New(`form_should_have_property_id`)
	}
	if data["type"] == nil {
		return nil, errors.New(`form_should_have_property_type`)
	}
	if data["title"] == nil {
		return nil, errors.New(`form_should_have_property_title`)
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

		err := validateInfo(info, typ.ValidInfoKeys)
		if err != nil {
			fmt.Println(data)
			return nil, err
		}
	}

	if typ.Type != "field" {
		return nil, errors.New(`field_type_not_supported`)
	}

	vals, err := extractValidations(validationInterface)
	if err != nil {
		return nil, err
	}

	var desc *string
	if data["description"] != nil {
		descVal := data["description"].(string)
		desc = &descVal
	}

	var placeholder *string
	if data["placeholder"] != nil {
		placeholderVal := data["placeholder"].(string)
		placeholder = &placeholderVal
	}

	return &Field{
		ID:          data["id"].(string),
		Type:        *typ,
		Title:       data["title"].(string),
		Description: desc,
		Placeholder: placeholder,
		Validations: vals,
		Info:        info,
	}, nil
}

func validateInfo(info map[string]interface{}, validInfoKeys []Info) error {
	for _, inf := range validInfoKeys {
		val, ok := info[inf.Key]
		if !ok && !inf.IsOptional {
			return errors.New(fmt.Sprintf("info should have property: %s", inf.Key))
		}

		if inf.FieldInfoValidation != nil {
			err := inf.FieldInfoValidation(val.(string))
			if nil != err {
				return err
			}
		}

		if len(inf.Child) == 0 {
			continue
		}

		valMap, ok := val.(map[string]interface{})
		if !ok {
			return errors.New(fmt.Sprintf("property should be an object: %s", inf.Key))
		}

		err := validateInfo(valMap, inf.Child)
		if nil != err {
			return err
		}
	}

	return nil
}

func extractValidations(vals []interface{}) ([]FieldValidation, error) {
	result := make([]FieldValidation, 0)
	for _, val := range vals {
		data, ok := val.(map[string]interface{})
		if !ok {
			return nil, errors.New("validation must be object")
		}

		v := FieldValidation{}

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