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

func (validator *Validator) Validate(structure map[string]interface{}) error {
	dittoType, ok := structure["type"]
	if !ok {
		return errors.New("type is expected on structure attribute")
	}

	err := validator.validateAttributes(dittoType, structure)
	if err != nil {
		return err
	}

	return nil
}

func (validator *Validator) validateAttributes(dittoType interface{}, structure map[string]interface{}) error {
	m, ok := validator.metadata[dittoType]
	if !ok {
		return errors.New(dittoType.(string) + " is not recognized in current ditto version: " + validator.version)
	}

	mapType, ok := m.(map[interface{}]interface{})
	typeAttributes, ok := mapType["attributes"]
	if !ok {
		return errors.New("attributes is expected in metadata: " + dittoType.(string))
	}

	mapTypeAttributes, ok := typeAttributes.(map[interface{}]interface{})
	if !ok {
		return errors.New("map attributes is expected in metadata: " + dittoType.(string))
	}

	for k, v := range mapTypeAttributes {
		e := validator.validateAttribute(k.(string), singleValueToArray(v), structure)
		if e != nil {
			return e
		}
	}
	return nil
}

func (validator *Validator) validateAttribute(key string, rules []interface{}, structure map[string]interface{}) error {
	_, ok := structure[key]
	if !ok {
		if !isIn(rules, "optional") {
			return errors.New("required: " + key)
		}
	} else {
		e := validator.checkRules(rules, structure, key)
		if e != nil {
			return e
		}
	}
	return nil
}

func (validator *Validator) checkRules(rules []interface{}, data map[string]interface{}, key string) error {
	for _, rule := range rules {
		ruleName := ""
		ruleDetails := make([]interface{}, 0)
		ruleName, isString := rule.(string)

		if !isString {
			m := rule.(map[interface{}]interface{})
			for k, v := range m {
				ruleName = k.(string)
				ruleDetails = singleValueToArray(v)
				break
			}
		}

		f, funcExist := rulesFunc[ruleName]
		if !funcExist {
			return errors.New("rule not found: " + ruleName)
		}

		err := f(validator, ruleDetails, data, key)
		if err != nil {
			return err
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
