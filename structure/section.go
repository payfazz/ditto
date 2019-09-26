package structure

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Section struct {
	ID           string                 `json:"id"`
	Type         Type                   `json:"type"`
	Title        string                 `json:"title"`
	Description  *string                `json:"description"`
	ChildSection []Section              `json:"child_section"`
	ChildField   []Field                `json:"child_field"`
	Info         map[string]interface{} `json:"info,omitempty"`
}

func (s Section) MarshalJSON() ([]byte, error) {
	type WithSection struct {
		ID           string                 `json:"id"`
		Type         Type                   `json:"type"`
		Title        string                 `json:"title"`
		Description  *string                `json:"description"`
		ChildSection []Section              `json:"child"`
		Info         map[string]interface{} `json:"info,omitempty"`
	}

	if len(s.ChildSection) > 0 {
		result := WithSection{
			ID:           s.ID,
			Type:         s.Type,
			Title:        s.Title,
			Description:  s.Description,
			ChildSection: s.ChildSection,
			Info:         s.Info,
		}

		return json.Marshal(result)
	}

	type WithField struct {
		ID          string                 `json:"id"`
		Type        Type                   `json:"type"`
		Title       string                 `json:"title"`
		Description *string                `json:"description"`
		ChildField  []Field                `json:"child"`
		Info        map[string]interface{} `json:"info,omitempty"`
	}

	result := WithField{
		ID:          s.ID,
		Type:        s.Type,
		Title:       s.Title,
		Description: s.Description,
		ChildField:  s.ChildField,
		Info:        s.Info,
	}

	return json.Marshal(result)
}

func NewSectionFromMap(data map[string]interface{}) (*Section, error) {
	if data["id"] == nil {
		return nil, errors.New(`section_should_have_property_id`)
	}
	if data["title"] == nil {
		return nil, errors.New(`section_should_have_property_title`)
	}
	if data["type"] == nil {
		return nil, errors.New(`section_should_have_property_type`)
	}
	if data["child"] == nil {
		return nil, errors.New(`section_should_have_property_child`)
	}

	fieldType, ok := data["type"].(string)
	if !ok {
		fmt.Println(data)
		return nil, errors.New(`section_type_not_supported`)
	}

	typ := GetType(fieldType)
	if nil == typ {
		fmt.Println(data)
		return nil, errors.New(`section_type_not_supported`)
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

	if typ.Type != "section" {
		fmt.Println(data)
		return nil, errors.New(`section_type_not_supported`)
	}

	childInterface, ok := data["child"].([]interface{})
	if !ok {
		return nil, errors.New("section_child_should_be_array")
	}

	var desc *string
	if data["description"] != nil {
		descVal := data["description"].(string)
		desc = &descVal
	}
	result := &Section{
		ID:           data["id"].(string),
		Type:         *typ,
		Title:        data["title"].(string),
		Description:  desc,
		ChildSection: nil,
		ChildField:   nil,
		Info:         info,
	}

	childSection := extractArrayMap(childInterface)

	if len(childSection) == 0 {
		return nil, errors.New(`section_should_have_child`)
	}

	firstChild := childSection[0]

	if firstChild["child"] == nil {
		//fields
		childs := make([]Field, 0)
		for _, child := range childSection {
			field, err := NewFieldFromMap(child)
			if err != nil {
				return nil, err
			}
			childs = append(childs, *field)
		}
		result.ChildField = childs
	} else {
		//sections
		childs := make([]Section, 0)
		for _, child := range childSection {
			field, err := NewSectionFromMap(child)
			if err != nil {
				return nil, err
			}
			childs = append(childs, *field)
		}
		result.ChildSection = childs
	}

	return result, nil
}
