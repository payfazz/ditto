package ditto

type Field struct {
	ID          string
	Type        string
	Title       string
	Description string
	Validations string
	Info        map[string]interface{}
}
