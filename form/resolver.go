package form

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	rError "github.com/payfazz/kitx/pkg/error"
)

//resolveSection for resolve section
func (formSvc *Service) resolveSection(ctx context.Context, structure map[string]interface{}) (map[string]interface{}, error) {
	childInterface, ok := structure["child"].([]interface{})
	if !ok {
		return nil, errors.New(`invalid_structure_section`)
	}

	childSection := formSvc.extractArrayMap(childInterface)
	firstChild := childSection[0]

	if firstChild["child"] == nil {
		for index, child := range childSection {
			nestedData, err := formSvc.resolveForm(ctx, child)
			if err != nil {
				return nil, err
			}
			childSection[index] = nestedData
		}
	} else {
		for index, child := range childSection {
			nestedData, err := formSvc.resolveSection(ctx, child)
			if err != nil {
				return nil, err
			}
			childSection[index] = nestedData
		}
	}

	structure["child"] = childSection
	return structure, nil
}

//resolveForm for resolve value form
func (formSvc *Service) resolveForm(ctx context.Context, formField map[string]interface{}) (map[string]interface{}, error) {
	objValue := make(map[string]interface{})
	formType, ok := formField["type"].(string)
	if !ok {
		return nil, rError.New(errors.New(`invalid_structure_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_form`)
	}
	if formField["info"] != nil {
		infoForm, ok := formField["info"].(map[string]interface{})
		if !ok {
			return nil, rError.New(errors.New(`invalid_structure_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_form`)
		}
		if infoForm["options"] != nil {
			infoOptions := infoForm["options"].(map[string]interface{})
			if !ok {
				return nil, rError.New(errors.New(`invalid_structure_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_form`)
			}
			infoType, ok := infoOptions["type"].(string)
			if !ok {
				return nil, rError.New(errors.New(`invalid_structure_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_form`)
			}
			infoValue, ok := infoOptions["value"].(string)
			if !ok {
				return nil, rError.New(errors.New(`invalid_structure_form`), rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_form`)
			}
			objValue["type"] = infoValue
			if infoType == "static" {
				if formType == "object_searchable_list" {
					objectModel, err := formSvc.repo.GetAllObjectListByCode(ctx, infoValue)
					if err != nil {
						return nil, err
					}
					if len(objectModel) == 0 {
						return nil, rError.New(errors.New(`list_option_`+infoValue+`_not_supported`), rError.Enum.UNPROCESSABLEENTITY, `list_option_`+infoValue+`_not_supported`)
					}
					arrObjVal := make([]map[string]interface{}, 0)
					for _, objval := range objectModel {
						arrObjVal = append(arrObjVal, map[string]interface{}{
							"id":    objval.ID,
							"value": objval.Name,
						})
					}
					objValue["value"] = arrObjVal
				} else {
					optionsModel, err := formSvc.repo.GetListOptionByType(ctx, infoValue)
					if err != nil {
						return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
					}
					if optionsModel == nil {
						return nil, rError.New(errors.New(`list_option_`+infoValue+`_not_supported`), rError.Enum.UNPROCESSABLEENTITY, `list_option_`+infoValue+`_not_supported`)
					}
					objValue["value"] = optionsModel.Value
				}
			} else {
				dynamicModel, err := formSvc.repo.GetDynamicOptionByCode(ctx, infoValue)
				if err != nil {
					return nil, rError.New(err, rError.Enum.INTERNALSERVERERROR, `something_went_wrong`)
				}
				if dynamicModel == nil {
					return nil, rError.New(errors.New(`list_option_`+infoValue+`_not_supported`), rError.Enum.UNPROCESSABLEENTITY, `list_option_`+infoValue+`_not_supported`)
				}
				var resultMetadata map[string]interface{}
				var requestBody map[string]interface{}
				err = json.Unmarshal([]byte(dynamicModel.Result), &resultMetadata)
				if err != nil {
					return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_jsonb_result`)
				}
				if dynamicModel.Request != "{}" {
					err = json.Unmarshal([]byte(dynamicModel.Request), &requestBody)
					if err != nil {
						return nil, rError.New(err, rError.Enum.UNPROCESSABLEENTITY, `invalid_structure_jsonb_request`)
					}
				}
				objValue["url"] = dynamicModel.URL
				objValue["result"] = resultMetadata
				objValue["request"] = requestBody
				objValue["method"] = dynamicModel.Method
			}
			infoOptions["value"] = objValue
			infoForm["options"] = infoOptions
		}
		formField["info"] = infoForm
	}
	return formField, nil
}

func (formSvc *Service) resolveValue(value, additionalData map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	for index, val := range value {
		key := matchFirstCap.ReplaceAllString(index, "${1}_${2}")
		key = strings.ToLower(matchAllCap.ReplaceAllString(key, "${1}_${2}"))
		value, _ := val.(map[string]interface{})
		formType, _ := value["type"].(string)
		formValue, _ := value["value"].(string)
		result[key] = formValue
		if formType == "object_searchable_list" {
			var objectValue map[string]interface{}
			_ = json.Unmarshal([]byte(formValue), &objectValue)
			result[key] = objectValue
			result[key+"_id"] = objectValue["id"]
			result[key+"_value"] = objectValue["value"]
		}
	}
	for index, _ := range additionalData {
		key := matchFirstCap.ReplaceAllString(index, "${1}_${2}")
		key = strings.ToLower(matchAllCap.ReplaceAllString(key, "${1}_${2}"))
		result[key] = additionalData[index]
	}
	return result
}
