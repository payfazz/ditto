package component

import (
	"errors"
	"fmt"
)

type Attributes map[string]interface{}

type Interface interface {
	RequiredAttrs() []string
	FillStruct(attrs Attributes) error
}

func ValidateRequiredAttrs(obj Interface, attrs Attributes) error {
	for _, req := range obj.RequiredAttrs() {
		if _, ok := attrs[req]; !ok {
			return errors.New(fmt.Sprintf("required attribute: %s", req))
		}
	}
	return nil
}

type Component struct{
	ID string
	Description string
	Title string
	Type string
}

func (c *Component) RequiredAttrs() []string {
	return []string{
		"id",
		"description",
		"title",
		"type",
	}
}

func (c *Component) FillStruct(attrs Attributes) error {
	if attrs["id"] == nil {
		return errors.New(`component should have property: id`)
	}

	c.ID = attrs["id"].(string)

	if attrs["type"] == nil {
		return errors.New(`component should have property: type`)
	}

	c.Type = attrs["type"].(string)

	if attrs["title"] == nil {
		return errors.New(`component should have property: title`)
	}

	c.Title = attrs["title"].(string)

	if attrs["description"] == nil {
		return errors.New(`component should have property: description`)
	}

	c.Description = attrs["description"].(string)

	return nil
}

