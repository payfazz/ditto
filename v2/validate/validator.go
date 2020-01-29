package validate

import (
	"github.com/go-yaml/yaml"
	ditto_yaml "github.com/payfazz/ditto-yaml"
	"github.com/payfazz/ditto/v2/validate/structure"
)

type Validator struct {
	metadata map[interface{}]interface{}
	version  string
	structValidator *structure.Validator
}

func New() *Validator {
	dittoVersion := "v0.1"

	metadata, err := ditto_yaml.Get(dittoVersion)
	if nil != err {
		panic(err)
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(metadata), &m)
	if nil != err {
		panic(err)
	}

	structValidator := structure.New(m, dittoVersion)

	return &Validator{
		metadata: m,
		version:  dittoVersion,
		structValidator: structValidator,
	}
}

func (v *Validator) ValidateStructure(structure map[string]interface{}) error {
	return v.structValidator.Validate(structure)
}