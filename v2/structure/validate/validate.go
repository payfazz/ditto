package validate

import (
	"errors"
	"github.com/go-yaml/yaml"
	ditto_yaml "github.com/payfazz/ditto-yaml"
)

type Validator struct {
	metadata map[interface{}]interface{}
	version  string
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

	return &Validator{
		metadata: m,
		version:  dittoVersion,
	}
}

func (validator *Validator) Validate(data map[string]interface{}) error {
	typ, ok := data["type"]
	if !ok {
		return errors.New("type is expected on data attribute")
	}

	m, ok := validator.metadata[typ]
	if !ok {
		return errors.New(typ.(string) + " is not recognized in current ditto version: " + validator.version)
	}

	mapType, ok := m.(map[interface{}]interface{})

	attr, ok := mapType["attributes"]
	if !ok {
		return errors.New("attributes is expected in metadata: " + typ.(string))
	}

	mapAttr, ok := attr.(map[interface{}]interface{})
	if !ok {
		return errors.New("map attributes is expected in metadata: " + typ.(string))
	}

	for k, v := range mapAttr {
		rules := singleValueToArray(v)
		_, ok := data[k.(string)]
		if !ok {
			if !isIn(rules, "optional") {
				return errors.New("required: " + k.(string))
			}
			continue
		}

		for _, rule := range rules {
			ruleName := ""
			ruleDetails := make([]interface{}, 0)
			ruleName, ok := rule.(string)

			if !ok {
				m := rule.(map[interface{}]interface{})
				for k, v := range m {
					ruleName = k.(string)
					ruleDetails = singleValueToArray(v)
					break
				}
			}

			f, ok := rulesFunc[ruleName]
			if !ok {
				return errors.New("rule not found: "  + ruleName)
			}

			err := f(validator, ruleDetails, data, k.(string))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isIn(haystack []interface{}, needle interface{}) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func singleValueToArray(val interface{}) []interface{} {
	if v, ok := val.([]interface{}); ok {
		return v
	}

	return []interface{}{val}
}
