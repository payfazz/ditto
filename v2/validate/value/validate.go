package value

import (
	"errors"
	"fmt"
)

type Validator struct {
	metadata map[interface{}]interface{}
	version  string
}

func New(metadata map[interface{}]interface{}, version string) *Validator {
	return &Validator{
		metadata: metadata,
		version:  version,
	}
}

func (validator *Validator) ValidateInput(inputs map[string]interface{}, fields map[string]interface{}) error {
	for id, field := range fields {
		validations, ok := field.([]interface{})
		if !ok {
			return errors.New("field is not an array: " + id)
		}

		for _, validation := range validations {
			mapValidation, ok := validation.(map[string]interface{})
			if !ok {
				return errors.New("validation is not a map: " + id)
			}

			typ, ok := mapValidation["type"].(string)
			if !ok {
				return errors.New("field does not have type: " + id)
			}

			rule := rules[typ]
			if rule == nil {
				continue
			}

			fmt.Println(typ, inputs)
			if !rule(inputs[id], mapValidation["value"]) {
				return errors.New(mapValidation["error_message"].(string) + " : " + id)
			}
		}

	}
	return nil
}

func (validator *Validator) ExtractField(structure map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	dittoType, ok := structure["type"]
	if !ok {
		return result, errors.New("type is expected on structure attribute")
	}

	attrs, err := validator.getMetadataByType(dittoType.(string))
	if err != nil {
		return result, err
	}

	descendants := validator.getDescendants(dittoType)
	if isIn(descendants, "field") {
		result[structure["id"].(string)] = structure["validations"]
	}

	for k, v := range attrs {
		rules := singleValueToArray(v)
		if validator.hasComposite(rules) {
			res, err := validator.extractFieldFromComposite(structure[k.(string)].([]interface{}))
			if nil != err {
				return result, err
			}

			for k, v := range res {
				result[k] = v
			}
		}
	}

	return result, nil
}

func (validator *Validator) hasComposite(rules []interface{}) bool {
	for _, rule := range rules {
		ruleName := ""
		ruleName, isString := rule.(string)

		if !isString {
			m := rule.(map[interface{}]interface{})
			for k, _ := range m {
				ruleName = k.(string)
				break
			}
		}

		if ruleName == "composite" {
			return true
		}
	}
	return false
}

func (validator *Validator) extractFieldFromComposite(composite []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, c := range composite {
		m, ok := c.(map[string]interface{})
		if !ok {
			return result, errors.New("map of interface is expected")
		}

		res, err := validator.ExtractField(m)
		if nil != err {
			return result, err
		}

		for k, v := range res {
			result[k] = v
		}
	}

	return result, nil
}

func (validator *Validator) getMetadataByType(dittoType string) (map[interface{}]interface{}, error) {
	m, ok := validator.metadata[dittoType]
	if !ok {
		return nil, errors.New(dittoType + " is not recognized in current ditto version: " + validator.version)
	}

	mapType, ok := m.(map[interface{}]interface{})
	typeAttributes, ok := mapType["attributes"]
	if !ok {
		return nil, errors.New("attributes is expected in metadata: " + dittoType)
	}

	mapTypeAttributes, ok := typeAttributes.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("map attributes is expected in metadata: " + dittoType)
	}

	return mapTypeAttributes, nil
}

func isIn(haystack []interface{}, needle interface{}) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func singleValueToArray(val interface{}) []interface{} {
	if v, ok := val.([]interface{}); ok {
		return v
	}

	return []interface{}{val}
}

func (validator *Validator) getDescendants(typ interface{}) []interface{} {
	if base, ok := validator.metadata[typ]; ok {
		m := base.(map[interface{}]interface{})
		parent := validator.getDescendants(m["base"])
		return append(parent, typ)
	}
	return []interface{}{}
}
