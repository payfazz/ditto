package form

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

//validateSection for validate form section
func (formSvc *Service) validateSection(ctx context.Context, data map[string]interface{}, usedIDs map[string]bool) error {
	var err error
	if data["id"] == nil {
		return errors.New(`section_should_have_property_id`)
	}
	_, ok := usedIDs[data["id"].(string)]
	if !ok {
		usedIDs[data["id"].(string)] = true
	} else {
		return errors.New(`field_id_must_unique`)
	}
	if data["title"] == nil {
		return errors.New(`section_should_have_property_title`)
	}
	if data["type"] == nil {
		return errors.New(`section_should_have_property_type`)
	}
	if data["child"] == nil {
		return errors.New(`section_should_have_property_child`)
	}

	typeSection, ok := data["type"].(string)
	if !ok {
		return errors.New(`section_type_should_be_string`)
	}
	sectionValid, err := formSvc.repo.GetSupportedTypeByTypeAndValue(ctx, "section", typeSection)
	if sectionValid == nil {
		if err != nil {
			return err
		}
		return errors.New(`section_type_not_supported`)
	}

	childInterface, ok := data["child"].([]interface{})
	if !ok {
		return errors.New("section_child_should_be_array")
	}

	childSection := formSvc.extractArrayMap(childInterface)

	firstChild := childSection[0]

	if firstChild["child"] == nil {
		var fieldID string
		for _, child := range childSection {
			err = formSvc.validateForm(ctx, child)
			if err != nil {
				return err
			}
			fieldID, ok = child["id"].(string)
			if !ok {
				return errors.New(`field_id_should_be_string`)
			}
			_, ok := usedIDs[fieldID]
			if !ok {
				usedIDs[fieldID] = true
			} else {
				return errors.New(`field_id_must_unique`)
			}
		}
	} else {
		for _, child := range childSection {
			err = formSvc.validateSection(ctx, child, usedIDs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//validateForm for validate form object
func (formSvc *Service) validateForm(ctx context.Context, data map[string]interface{}) error {
	var err error
	if data["id"] == nil {
		return errors.New(`form_should_have_property_id`)
	}
	if data["type"] == nil {
		return errors.New(`form_should_have_property_type`)
	}
	if data["title"] == nil {
		return errors.New(`form_should_have_property_name`)
	}
	if data["description"] == nil {
		return errors.New(`form_should_have_property_description`)
	}
	if data["validations"] == nil {
		return errors.New(`form_should_have_property_validation`)
	}

	formType, ok := data["type"].(string)
	if !ok {
		return errors.New(`form_type_should_be_string`)
	}
	formValid, err := formSvc.repo.GetSupportedTypeByTypeAndValue(ctx, "form", formType)
	if formValid == nil {
		if err != nil {
			return err
		}
		return errors.New(`form_type_not_supported`)
	}

	if data["info"] != nil {
		infoForm, ok := data["info"].(map[string]interface{})
		if !ok {
			return errors.New("form_info_should_be_an_object")
		}
		err = formSvc.validateInfoForm(ctx, infoForm, formType)
		if err != nil {
			return err
		}
	}

	validationInterface, ok := data["validations"].([]interface{})
	if !ok {
		return errors.New("form_validation_should_be_array")
	}
	validationRules := formSvc.extractArrayMap(validationInterface)

	for _, rule := range validationRules {
		err := formSvc.validateFormValidation(rule)
		if err != nil {
			return err
		}
	}

	return nil
}

//validateInforForm for validate form info object
func (formSvc *Service) validateInfoForm(ctx context.Context, info map[string]interface{}, formType string) error {
	if formType == "photo_camera" {
		// if info["instruction_text"] == nil {
		// 	return errors.New("form_info_photo_camera_should_have_property_instruction_text")
		// }
		if info["instruction_image"] == nil {
			return errors.New("form_info_photo_camera_should_have_property_instruction_image")
		}
	}
	result, err := formSvc.repo.GetSupportedTypeByValue(ctx, formType)
	if err != nil {
		return err
	}
	if result.Group == "list" {
		if info["options"] == nil {
			return errors.New("form_info_list_should_have_property_options")
		}
		infoOptions, ok := info["options"].(map[string]interface{})
		if !ok {
			return errors.New("form_info_options_should_be_an_object")
		}
		if infoOptions["type"] == nil {
			return errors.New("form_info_options_should_have_property_type")
		}
		typeOptions, ok := infoOptions["type"].(string)
		if !ok {
			return errors.New(`info_type_should_be_string`)
		}
		if typeOptions != "static" && typeOptions != "dynamic" {
			return errors.New("form_info_options_type_should_be_static_or_dynamic")
		}
		if infoOptions["value"] == nil {
			return errors.New("form_info_options_should_have_property_value")
		}
		// valueInfo, ok := infoOptions["value"].(map[string]interface{})
		// if !ok {
		// 	return errors.New("form_info_options_value_should_be_an_object")
		// }
		// if valueInfo["type"] == nil {
		// 	return errors.New("form_info_options_value_should_have_property_type")
		// }
		// if valueInfo["value"] == nil {
		// 	return errors.New("form_info_options_value_should_have_property_value")
		// }
		if typeOptions == "static" {
			valueOptions, ok := infoOptions["value"].(string)
			if !ok {
				return errors.New(`info_options_value_should_be_string`)
			}
			if formType == "object_searchable_list" {
				resultObjList, _ := formSvc.repo.GetObjectListByCode(ctx, valueOptions)
				if resultObjList == nil {
					return fmt.Errorf(`form_info_options_value_%s_not_supported`, valueOptions)
				}
			} else {
				listValid, _ := formSvc.repo.GetListOptionByType(ctx, valueOptions)
				if listValid == nil {
					return fmt.Errorf(`form_info_options_value_%s_not_supported`, valueOptions)
				}
			}
		} else {
			valueOptions, ok := infoOptions["value"].(string)
			if !ok {
				return errors.New(`info_options_value_should_be_string`)
			}
			dynamicValid, _ := formSvc.repo.GetDynamicOptionByCode(ctx, valueOptions)
			if dynamicValid == nil {
				return fmt.Errorf(`form_info_options_value_%s_not_supported`, valueOptions)
			}
		}
	}
	return nil
}

//validateFormValidation for validate supported validation type
func (formSvc *Service) validateFormValidation(validation map[string]interface{}) error {
	if validation["type"] == nil {
		return errors.New(`form_validation_should_have_property_type`)
	}
	if validation["error_message"] == nil {
		return errors.New(`form_validation_should_have_property_error_message`)
	}
	ruleType, ok := validation["type"].(string)
	if !ok {
		return errors.New(`rule_type_should_be_string`)
	}
	if ruleType != "required" {
		if validation["value"] == nil {
			return errors.New(`form_validation_should_have_property_value`)
		}
		if strings.Contains(ruleType, "between") {
			ruleValue, ok := validation["value"].(string)
			if !ok {
				return errors.New(`rule_value_should_be_string`)
			}
			splitedString := strings.Split(ruleValue, ",")
			if len(splitedString) != 2 {
				return errors.New(`form_validation_between_value_should_splited_with_comma`)
			}
		}
	}
	return nil
}
