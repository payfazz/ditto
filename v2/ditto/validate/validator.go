package validate

import (
	"github.com/go-yaml/yaml"
	ditto_yaml "github.com/payfazz/ditto-yaml"
	structure2 "github.com/payfazz/ditto/v2/ditto/validate/structure"
	value2 "github.com/payfazz/ditto/v2/ditto/validate/value"
)

type Validator struct {
	metadata        map[interface{}]interface{}
	version         string
	structValidator *structure2.Validator
	valueValidator  *value2.Validator
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

	structValidator := structure2.New(m, dittoVersion)
	valueValidator := value2.New(m, dittoVersion)

	return &Validator{
		metadata:        m,
		version:         dittoVersion,
		structValidator: structValidator,
		valueValidator:  valueValidator,
	}
}

func (v *Validator) ValidateStructure(structure map[string]interface{}) error {
	return v.structValidator.Validate(structure)
}

func (v *Validator) ValidateInput(inputs map[string]interface{}, structure map[string]interface{}) error {
	fields, err := v.valueValidator.ExtractField(structure)
	if nil != err {
		return err
	}
	return v.valueValidator.ValidateInput(inputs, fields)
}

func (v *Validator) ExtractField(structure map[string]interface{}) (map[string]interface{}, error) {
	return v.valueValidator.ExtractField(structure)
}
