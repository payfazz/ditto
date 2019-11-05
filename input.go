package ditto

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ValidatedInput struct {
	IsValid   bool                   `json:"valid"`
	Structure string                 `json:"structure"`
	Values    map[string]interface{} `json:"values"`
	Errors    map[string]string      `json:"field_error"`
}

func ValidateFormInput(root *Section, input map[string]interface{}) (*ValidatedInput, error) {
	validatedRoot := &SectionWithStatus{
		ID:          root.ID,
		Type:        root.Type,
		Title:       root.Title,
		Description: root.Description,
		Info:        root.Info,
	}
	errMap := make(map[string]string)
	errVal := validateFormInput(root, input, validatedRoot, errMap)

	structureJson, err := json.Marshal(validatedRoot)
	if nil != err {
		return nil, err
	}

	return &ValidatedInput{
		IsValid:   nil == errVal,
		Structure: string(structureJson),
		Values:    input,
		Errors:    errMap,
	}, nil
}

func validateFormInput(root *Section, input map[string]interface{}, validatedRoot *SectionWithStatus, errMap map[string]string) error {
	if len(root.ChildSection) > 0 {
		errs := make([]string, 0)
		childs := make([]*SectionWithStatus, 0)
		for _, child := range root.ChildSection {
			sec := &SectionWithStatus{
				ID:          child.ID,
				Type:        child.Type,
				Title:       child.Title,
				Description: child.Description,
				Info:        child.Info,
			}
			childs = append(childs, sec)
			_validatedRoot := &SectionWithStatus{}
			err := validateFormInput(&child, input, _validatedRoot, errMap)
			if nil != err {
				errs = append(errs, err.Error())
				sec.Status = map[string]interface{}{
					"type":    "error",
					"message": err.Error(),
				}
			}
		}

		validatedRoot.ChildSection = childs
		if len(errs) > 0 {
			errStr := strings.Join(errs, ",")
			validatedRoot.Status = map[string]interface{}{
				"type":    "error",
				"message": errStr,
			}

			return errors.New(errStr)
		}

		return nil
	}

	errs := make([]string, 0)
	childs := make([]*FieldWithValue, 0)
	for _, child := range root.ChildField {
		fi := &FieldWithValue{
			ID:          child.ID,
			Type:        child.Type,
			Title:       child.Title,
			Description: child.Description,
			Placeholder: child.Placeholder,
			Validations: child.Validations,
			Info:        child.Info,
			Status:      child.Status,
		}

		childs = append(childs, fi)

		values, ok := input[child.ID]
		if !ok {
			err := errors.New(fmt.Sprintf("input not found: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			fi.Status = errObj
			errMap[child.ID] = err.Error()
			errs = append(errs, err.Error())
			continue
		}

		valueMap, ok := values.(map[string]interface{})
		if !ok {
			err := errors.New(fmt.Sprintf("input must be an object: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			fi.Status = errObj
			errMap[child.ID] = err.Error()
			errs = append(errs, err.Error())
			continue
		}

		value := valueMap["value"]
		if nil == value {
			err := errors.New(fmt.Sprintf("value property not found: %s", child.ID))
			errObj := make(map[string]interface{})
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			fi.Status = errObj
			errMap[child.ID] = err.Error()
			errs = append(errs, err.Error())
			continue
		}

		fi.Value = value

		for _, validation := range child.Validations {
			validator := GetValidator(validation.Type)
			ok := validator(valueMap, validation)
			if !ok {
				errObj := make(map[string]interface{})
				errObj["type"] = "error"
				errObj["message"] = validation.ErrorMessage
				valueMap["error"] = validation.ErrorMessage
				fi.Status = errObj
				errMap[child.ID] = validation.ErrorMessage
				errs = append(errs, validation.ErrorMessage)
				break
			}
		}
	}

	validatedRoot.ChildField = childs
	if len(errs) > 0 {
		errStr := strings.Join(errs, ",")
		validatedRoot.Status = map[string]interface{}{
			"type":    "error",
			"message": errStr,
		}

		return errors.New(errStr)
	}

	return nil
}
