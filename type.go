package ditto

import "errors"

var ErrTypeExists = errors.New("type already exists")

type Type struct {
	Type  string
	Value string
	Group string
}

var types map[string]Type

func RegisterType(t Type) error {
	if _, ok := types[t.Value]; ok {
		return ErrTypeExists
	}

	types[t.Value] = t
	return nil
}
