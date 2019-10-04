package component

type Interface interface {
	RequiredAttrs() []string
}

type Component struct{}

func (c *Component) RequiredAttrs() []string {
	return []string{
		"id",
		"description",
		"title",
		"type",
	}
}
