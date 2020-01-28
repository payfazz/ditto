package validate

import (
	"errors"
)

type RuleFunc func(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error

var rulesFunc map[string]RuleFunc

func init() {
	rulesFunc = make(map[string]RuleFunc)

	rulesFunc["composite"] = composite
	rulesFunc["required"] = alwaysNil
	rulesFunc["optional"] = alwaysNil
}

func alwaysNil(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error {
	return nil
}

func composite(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error {
	val := data[key]
	vals, ok := val.([]interface{})
	if !ok {
		return errors.New("array of interface is expected in composite")
	}

	for _, v := range vals {
		m, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("map is expected in composite")
		}

		typ, ok := m["type"]
		if !ok {
			return errors.New("type is expected on data attribute")
		}

		descendants := validator.getDescendants(typ)
		if !isIn(descendants, typ) {
			return errors.New(typ.(string) + " is not recognized in current ditto version: " + validator.version)
		}
	}

	return nil
}

func (validator *Validator) getDescendants(typ interface{}) []interface{} {
	if base, ok := validator.metadata[typ]; ok {
		m := base.(map[interface{}]interface{})
		parent := validator.getDescendants(m["base"])
		return append(parent, typ)
	}
	return []interface{}{}
}
