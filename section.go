package ditto

import (
	"errors"
	"fmt"
)

type Section struct {
	ID           string                 `json:"id"`
	Type         Type                   `json:"type"`
	Title        string                 `json:"title"`
	ChildSection []Section              `json:"child_section"`
	ChildField   []Field                `json:"child_field"`
	Info         map[string]interface{} `json:"info"`
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

	result := &Section{
		ID:           data["id"].(string),
		Type:         *typ,
		Title:        data["title"].(string),
		ChildSection: nil,
		ChildField:   nil,
		Info:         info,
	}

	childSection := extractArrayMap(childInterface)

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
