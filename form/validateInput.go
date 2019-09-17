package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

//validateAndGenerateFormInput for validateAndGenerate value form section
func (formSvc *Service) validateAndGenerateFormInput(ctx context.Context, data map[string]interface{}, structure map[string]interface{}, fieldError *[]map[string]string) (map[string]interface{}, error) {
	var errorString []string
	childInterface, ok := structure["child"].([]interface{})
	if !ok {
		return nil, errors.New("section_child_should_be_array")
	}

	childSection := formSvc.extractArrayMap(childInterface)
	firstChild := childSection[0]

	if firstChild["child"] == nil {
		for index, child := range childSection {
			formFieldID, _ := child["id"].(string)
			nestedData, err := formSvc.validateAndGenerateFormWithValidations(ctx, data, child)
			if err != nil {
				errorString = append(errorString, err.Error())
				*fieldError = append(*fieldError, map[string]string{
					formFieldID: err.Error(),
				})
			}
			childSection[index] = nestedData
		}
	} else {
		for index, child := range childSection {
			nestedData, err := formSvc.validateAndGenerateFormInput(ctx, data, child, fieldError)
			if err != nil {
				errorString = append(errorString, err.Error())
			}
			childSection[index] = nestedData
		}
	}
	var childError error
	if len(errorString) > 0 {
		errObj := make(map[string]interface{})
		errObj["type"] = "error"
		errObj["message"] = strings.Join(errorString, ", ")
		structure["status"] = errObj
		childError = errors.New(errObj["message"].(string))
	}

	structure["child"] = childSection
	return structure, childError
}

//validateAndGenerateFormWithValidations for validateAndGenerate value with validation rules
func (formSvc *Service) validateAndGenerateFormWithValidations(ctx context.Context, formValues map[string]interface{}, formField map[string]interface{}) (map[string]interface{}, error) {
	var err error
	errObj := make(map[string]interface{})
	formField, err = formSvc.resolveForm(ctx, formField)
	if err != nil {
		errObj["type"] = "error"
		errObj["message"] = err.Error()
		formField["status"] = errObj
		return formField, err
	}

	formFieldID, ok := formField["id"].(string)
	if !ok {
		errObj["type"] = "error"
		errObj["message"] = "form_field_id_should_be_string"
		formField["status"] = errObj
		return formField, errors.New("form_field_id_should_be_string")
	}

	typeField, ok := formField["type"].(string)
	if !ok {
		errObj["type"] = "error"
		errObj["message"] = "form_field_type_should_be_string"
		formField["status"] = errObj
		return formField, errors.New("form_field_type_should_be_string")
	}

	result, err := formSvc.repo.GetSupportedTypeByTypeAndValue(ctx, "form", typeField)
	if err != nil {
		errObj["type"] = "error"
		errObj["message"] = err.Error()
		formField["status"] = errObj
		return formField, err
	}

	if result == nil {
		errObj["type"] = "error"
		errObj["message"] = errors.New(`form_type_not_supported`)
		formField["status"] = errObj
		return formField, err
	}

	currentValue, ok := formValues[formFieldID].(map[string]interface{})
	if ok {
		if currentValue["status"] != nil {
			statusValue, ok := currentValue["status"].(string)
			if !ok {
				errObj["type"] = "error"
				errObj["message"] = "status_should_be_string"
				formField["status"] = errObj
				return formField, errors.New("status_should_be_string")
			}
			if statusValue == "error" {
				errObj["type"] = "error"
				errObj["message"] = currentValue["error_message"].(string)
				formField["status"] = errObj
				return formField, errors.New(currentValue["error_message"].(string))
			}
		}
	}

	currentType, ok := currentValue["type"].(string)
	if !ok {
		errObj["type"] = "error"
		errObj["message"] = "type_should_be_string"
		formField["status"] = errObj
		return formField, errors.New("type_should_be_string")
	}

	if typeField != currentType {
		fmt.Println(typeField, currentType)
		errObj["type"] = "error"
		errObj["message"] = "type_not_match"
		formField["status"] = errObj
		return formField, errors.New("type_not_match")
	}

	validationInterface, ok := formField["validations"].([]interface{})
	if !ok {
		errObj["type"] = "error"
		errObj["message"] = "form_validation_should_be_array"
		formField["status"] = errObj
		return formField, errors.New("form_validation_should_be_array")
	}
	validationRules := formSvc.extractArrayMap(validationInterface)

	for _, rule := range validationRules {
		err := formSvc.validateValues(rule["type"].(string), rule["value"], formValues[formFieldID])
		if err != nil {
			errObj["type"] = "error"
			errObj["message"] = rule["error_message"].(string)
			formField["status"] = errObj
			return formField, errors.New(rule["error_message"].(string))
		}
		valueField := formValues[formFieldID].(map[string]interface{})
		formField["value"] = valueField["value"]
	}

	if result.Group == "list" {
		valueField, ok := formValues[formFieldID].(map[string]interface{})
		if !ok {
			errObj["type"] = "error"
			errObj["message"] = "form_value_should_be_json"
			formField["status"] = errObj
			return formField, errors.New("form_value_should_be_json")
		}
		formInfo, ok := formField["info"].(map[string]interface{})
		if !ok {
			errObj["type"] = "error"
			errObj["message"] = "form_info_should_be_json"
			formField["status"] = errObj
			return formField, errors.New("form_info_should_be_json")
		}
		formOptions, ok := formInfo["options"].(map[string]interface{})
		if !ok {
			errObj["type"] = "error"
			errObj["message"] = "form_info_options_should_be_json"
			formField["status"] = errObj
			return formField, errors.New("form_info_options_should_be_json")
		}
		formOptionsType := formOptions["type"].(string)
		typeList, ok := valueField["type_list"].(string)
		if !ok {
			typeList, ok = valueField["typeList"].(string)
			if !ok {
				errObj["type"] = "error"
				errObj["message"] = "type_list_not_recognized"
				formField["status"] = errObj
				return formField, errors.New("type_list_not_recognized")
			}
		}
		if typeList != formOptionsType {
			errObj["type"] = "error"
			errObj["message"] = "option_type_not_match"
			formField["status"] = errObj
			return formField, errors.New("option_type_not_match")
		}
		formOptionsValue, ok := formOptions["value"].(map[string]interface{})
		if !ok {
			errObj["type"] = "error"
			errObj["message"] = "form_info_options_value_should_be_json"
			formField["status"] = errObj
			return formField, errors.New("form_info_options_value_should_be_json")
		}
		formOptionsValueType := formOptionsValue["type"].(string)
		keyList, ok := valueField["key"].(string)
		if !ok {
			errObj["type"] = "error"
			errObj["message"] = "list_not_recognized"
			formField["status"] = errObj
			return formField, errors.New("list_not_recognized")
		}

		if formOptionsValueType != keyList {
			errObj["type"] = "error"
			errObj["message"] = "option_key_not_match"
			formField["status"] = errObj
			return formField, errors.New("option_key_not_match")
		}
		err = formSvc.validateListValue(ctx, formFieldID, valueField)
		if err != nil {
			errObj["type"] = "error"
			errObj["message"] = err.Error()
			formField["status"] = errObj
			return formField, err
		}
	}

	return formField, nil
}

func (formSvc *Service) validateListValue(ctx context.Context, idField string, valueField map[string]interface{}) error {
	found := false
	typeList, ok := valueField["type_list"].(string)
	if !ok {
		typeList, ok = valueField["typeList"].(string)
		if !ok {
			return errors.New("type_list_not_recognized")
		}
	}
	keyList, ok := valueField["key"].(string)
	if !ok {
		return errors.New(`list_not_recognized`)
	}

	formType, ok := valueField["type"].(string)
	if !ok {
		return errors.New(`form_not_recognized`)
	}

	valueList, ok := valueField["value"].(string)
	if !ok {
		return errors.New("value_list_not_valid")
	}
	if typeList == "static" {
		if formType == "object_searchable_list" {
			var objectValue map[string]interface{}
			err := json.Unmarshal([]byte(valueList), &objectValue)
			if err != nil {
				return errors.New("value_list_not_valid")
			}
			keyValueList, ok := objectValue["id"].(string)
			if !ok {
				return errors.New("value_list_not_valid")
			}
			uuidVal, err := uuid.Parse(keyValueList)
			if err != nil {
				return errors.New(`value_+` + keyList + `_not_valid`)
			}
			resultObjList, err := formSvc.repo.GetObjectListByID(ctx, uuidVal)
			if err != nil {
				return err
			}
			if resultObjList != nil {
				found = true
			}
		} else {
			resultList, err := formSvc.repo.GetListOptionByType(ctx, keyList)
			if err != nil {
				return err
			}
			if resultList == nil {
				return errors.New(`list_not_recognized`)
			}
			for _, val := range resultList.Value {
				if val == valueList {
					found = true
				}
			}
		}
	} else {
		resultList, err := formSvc.repo.GetDynamicOptionByCode(ctx, keyList)
		if err != nil {
			return err
		}
		if resultList == nil {
			return errors.New(`list_not_recognized`)
		}
		if keyList == "area" {
			var objectValue map[string]interface{}
			err := json.Unmarshal([]byte(valueList), &objectValue)
			if err != nil {
				return errors.New("value_list_not_valid")
			}
			valueList, ok = objectValue["id"].(string)
			if !ok {
				return errors.New("value_list_not_valid")
			}
		}
		resultValidator, err := formSvc.validate.Validate(ctx, resultList.URLValidate, valueList)
		if err != nil {
			return err
		}
		if resultValidator.Valid {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("invalid_value_for_%s", idField)
	}
	return nil
}

//validateValues for validate and generate value by type validation
func (formSvc *Service) validateValues(typeValidation string, validationValue interface{}, value interface{}) error {
	if typeValidation == "required" {
		if value == nil {
			return errors.New("error_required")
		}
		valueObj, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("value_should_be_json")
		}
		if valueObj["value"] == nil {
			return errors.New("error_required")
		}
	}
	valueObj := value.(map[string]interface{})
	if strings.Contains(typeValidation, "between") {
		validationVal := validationValue.(string)
		splitBetween := strings.Split(validationVal, ",")
		min, _ := strconv.Atoi(splitBetween[0])
		max, _ := strconv.Atoi(splitBetween[1])

		valueString, ok := valueObj["value"].(string)
		if !ok {
			return errors.New("value_list_not_valid")
		}

		if typeValidation == "text_between" {
			if len(valueString) < min || len(valueString) > max {
				return errors.New("error_text_between")
			}
		}
		if typeValidation == "count_between" {
			if len(valueString) < min || len(valueString) > max {
				return errors.New("error_count_between")
			}
		}
		if typeValidation == "age_between" {
			birthDate, _ := time.Parse("02-01-2006", valueString)
			age := formSvc.getAge(birthDate)
			if age < min || age > max {
				return errors.New("error_age_between")
			}
		}
		if typeValidation == "date_between" {
			birthDate, _ := time.Parse("02-01-2006", valueString)
			min, _ := time.Parse("02-01-2006", splitBetween[0])
			max, _ := time.Parse("02-01-2006", splitBetween[1])
			if birthDate.UTC().Unix() < min.UTC().Unix() || birthDate.UTC().Unix() > max.UTC().Unix() {
				fmt.Print(typeValidation+"\n", validationVal, "\n", birthDate, "\n")
				return errors.New("error_date_between")
			}
		}
	}
	if typeValidation == "regex" {
		validationVal, ok := validationValue.(string)
		if !ok {
			return errors.New("regex_value_should_be_string")
		}
		valueString, ok := valueObj["value"].(string)
		if !ok {
			return errors.New("value_should_be_string")
		}
		validationVal, _ = url.QueryUnescape(validationVal)
		match, _ := regexp.MatchString(validationVal, valueString)
		if !match {
			return errors.New("regexp_not_match")
		}
	}
	return nil
}
