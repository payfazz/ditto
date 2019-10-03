package component

type Interface interface {
	RequiredKeys() []string
	ComponentGroup() string
}

type Component struct{}

func (c *Component) RequiredKeys() []string {
	return []string{
		"id",
		"description",
		"title",
		"type",
	}
}
