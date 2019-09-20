package ditto

import (
	"errors"
)

var ErrTypeExists = errors.New("type already exists")
var ErrTypeCannotBeEmpty = errors.New("type cannot be empty")
var ErrValueCannotBeEmpty = errors.New("value cannot be empty")
var ErrGroupNotRegistered = errors.New("group is not registered")

type Info struct {
	Key            string
	Child          []Info
	InfoValidation func(data interface{}) error
	IsOptional     bool
}

type Type struct {
	Type  string
	Value string
	Group *Group
	Infos []Info
}

type Group struct {
	Name  string
	Infos []Info
}

var groups = map[string]*Group{}

func RegisterGroup(g *Group) error {
	if _, ok := groups[g.Name]; ok {
		return ErrTypeExists
	}

	groups[g.Name] = g
	return nil
}

func GetGroup(key string) *Group {
	return groups[key]
}

var types = map[string]*Type{}

func GetType(key string) *Type {
	return types[key]
}

func RegisterType(t *Type) error {
	if t.Value == "" {
		return ErrValueCannotBeEmpty
	}

	if t.Type == "" {
		return ErrTypeCannotBeEmpty
	}

	if _, ok := types[t.Value]; ok {
		return ErrTypeExists
	}

	if t.Group == nil {
		return ErrGroupNotRegistered
	}

	if GetGroup(t.Group.Name) == nil {
		return ErrGroupNotRegistered
	}

	types[t.Value] = t
	return nil
}

func init() {
	registerDefaultGroups()
	registerDefaultTypes()
}

func registerDefaultGroups() {
	_ = RegisterGroup(&Group{
		Name:  "section",
		Infos: nil,
	})

	_ = RegisterGroup(&Group{
		Name:  "section_field",
		Infos: nil,
	})

	_ = RegisterGroup(&Group{
		Name:  "text",
		Infos: nil,
	})

	_ = RegisterGroup(&Group{
		Name:  "file",
		Infos: nil,
	})

	_ = RegisterGroup(&Group{
		Name: "list",
		Infos: []Info{
			{
				Key: "options",
				Child: []Info{
					{
						Key: "type",
						InfoValidation: func(data interface{}) error {
							val, ok := data.(string)
							if !ok {
								return errors.New("info_type_should_be_string")
							}
							if val != "static" && val != "dynamic" {
								return errors.New("form_info_options_type_should_be_static_or_dynamic")
							}
							return nil
						},
					},
					{
						Key: "value",
					},
				},
			},
		},
	})
}

func registerDefaultTypes() {
	_ = RegisterType(&Type{
		Type:  "section",
		Value: "summary_field",
		Group: GetGroup("section_field"),
	})

	_ = RegisterType(&Type{
		Type:  "section",
		Value: "nextable_section",
		Group: GetGroup("section"),
	})

	_ = RegisterType(&Type{
		Type:  "section",
		Value: "nextable_field",
		Group: GetGroup("section_field"),
	})

	_ = RegisterType(&Type{
		Type:  "section",
		Value: "summary_section_send",
		Group: GetGroup("section_field"),
	})

	_ = RegisterType(&Type{
		Type:  "section",
		Value: "summary_section_save",
		Group: GetGroup("section_field"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "text_multiline",
		Group: GetGroup("text"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "text_numeric",
		Group: GetGroup("text"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "photo_camera",
		Group: GetGroup("file"),
		Infos: []Info{
			{
				Key:            "instruction_image",
				Child:          nil,
				InfoValidation: nil,
			},
		},
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "list",
		Group: GetGroup("list"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "text",
		Group: GetGroup("text"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "date",
		Group: GetGroup("text"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "searchable_list",
		Group: GetGroup("list"),
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "object_searchable_list",
		Group: GetGroup("list"),
		Infos: []Info{
			{
				Key: "type",
				InfoValidation: func(data interface{}) error {
					val, ok := data.(string)
					if !ok {
						return errors.New("info_type_should_be_string")
					}
					if val != "static" && val != "dynamic" {
						return errors.New("form_info_options_type_should_be_static_or_dynamic")
					}
					return nil
				},
			},
			{
				Key: "value",
			},
		},
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "normal_list",
		Group: GetGroup("list"),
	})

}
