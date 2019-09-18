package ditto

import "errors"

var ErrTypeExists = errors.New("type already exists")
var ErrGroupNotRegistered = errors.New("group is not registered")

type Type struct {
	Type      string
	Value     string
	Group     *Group
	ValidInfo []string
}

type Group struct {
	Name      string
	ValidInfo []string
}

var groups map[string]*Group

func RegisterGroup(g *Group) error {
	if _, ok := types[g.Name]; ok {
		return ErrTypeExists
	}

	types[g.Name] = g
	return nil
}

func GetGroup(key string) *Group {
	return groups[key]
}

var types map[string]*Type

func GetType(key string) *Type {
	return types[key]
}

func RegisterType(t *Type) error {
	if _, ok := types[t.Type]; ok {
		return ErrTypeExists
	}

	if t.Group == nil {
		return ErrGroupNotRegistered
	}

	if GetGroup(t.Group.Name) == nil {
		return ErrGroupNotRegistered
	}

	t.ValidInfo = append(t.ValidInfo, t.Group.ValidInfo...)
	types[t.Value] = t
	return nil
}

func init() {
	registerDefaultGroups()
	registerDefaultTypes()
}

func registerDefaultGroups() {
	_ = RegisterGroup(&Group{
		Name:      "section",
		ValidInfo: nil,
	})

	_ = RegisterGroup(&Group{
		Name:      "section_field",
		ValidInfo: nil,
	})

	_ = RegisterGroup(&Group{
		Name:      "text",
		ValidInfo: nil,
	})

	_ = RegisterGroup(&Group{
		Name:      "file",
		ValidInfo: nil,
	})

	_ = RegisterGroup(&Group{
		Name:      "list",
		ValidInfo: nil,
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
		Type:      "field",
		Value:     "photo_camera",
		Group:     GetGroup("file"),
		ValidInfo: []string{"instruction_image"},
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
	})

	_ = RegisterType(&Type{
		Type:  "field",
		Value: "normal_list",
		Group: GetGroup("list"),
	})

}
