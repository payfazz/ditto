package validate

import (
	"errors"
)

type RuleFunc func(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error

var rulesFunc map[string]RuleFunc

func init() {
	rulesFunc = make(map[string]RuleFunc)

	rulesFunc["required"] = alwaysNil
	rulesFunc["optional"] = alwaysNil

	rulesFunc["composite"] = composite
	rulesFunc["list"] = list
}

func alwaysNil(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error {
	return nil
}

func list(validator *Validator, ruleDetails []interface{}, data map[string]interface{}, key string) error {
	val := data[key]
	vals, ok := val.([]interface{})
	if !ok {
		return errors.New("array of interface is expected in composite")
	}

	t := ruleDetails[0].(string)

	for _, v := range vals {
		m, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("map of interface is expected")
		}

		err := validator.validateAttributes(t, m)
		if nil != err {
			return err
		}
	}

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

		err := validator.Validate(m)
		if nil != err {
			return err
		}

		typ, ok := m["type"]
		if !ok {
			return errors.New("type is expected on data attribute")
		}

		descendants := validator.getDescendants(typ)

		in := false
		for _, det := range ruleDetails {
			if isIn(descendants, det) {
				in = true
				break
			}
		}

		if !in {
			return errors.New(typ.(string) + " is not allowed in composite")
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
