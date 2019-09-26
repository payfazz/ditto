package structure_test

import (
	"github.com/payfazz/ditto/structure"
	"testing"
)

func TestRegisterFieldFail(t *testing.T) {
	typ := &structure.Type{}
	err := structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &structure.Type{
		Type: "test",
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &structure.Type{
		Value: "test",
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestRegisterGroupAndField(t *testing.T) {
	g := &structure.Group{
		Name:          "test",
		ValidInfoKeys: nil,
	}
	err := structure.RegisterGroup(g)
	if err != nil {
		t.Fatal(err)
	}

	err = structure.RegisterGroup(g)
	if err == nil {
		t.Fatal("error expected")
	}

	typ := &structure.Type{
		Type:          "test",
		Value:         "empty",
		Group:         g,
		ValidInfoKeys: nil,
	}
	err = structure.RegisterType(typ)
	if err != nil {
		t.Fatal(err)
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}
