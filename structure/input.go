package structure

import (
	"errors"
	"fmt"
)

func ValidateFormInput(root *Section, input map[string]interface{}) {
	if len(root.ChildSection) > 0 {
		for _, child := range root.ChildSection {
			ValidateFormInput(&child, input)
		}
		return
	}

	for _, child := range root.ChildField {
		values, ok := input[child.ID]
		if !ok {
			err := errors.New(fmt.Sprintf("input not found: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			child.Status = errObj
			continue
		}

		valueMap, ok := values.(map[string]interface{})
		if !ok {
			err := errors.New(fmt.Sprintf("input must be an object: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			child.Status = errObj
			continue
		}

		value := valueMap["value"]
		if nil == value {
			err := errors.New(fmt.Sprintf("value property not found: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			child.Status = errObj
			continue
		}

		for _, validation := range child.Validations {
			validator := GetValidator(validation.Type)
			ok := validator(valueMap, validation)
			if !ok {
				errObj := make(map[string]interface{})
				errObj["type"] = "error"
				errObj["message"] = validation.ErrorMessage
				valueMap["error"] = validation.ErrorMessage
				child.Status = errObj
			}
		}
	}
}
