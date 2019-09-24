package ditto_test

import (
	"github.com/payfazz/ditto"
	"testing"
)

func TestRegisterFieldFail(t *testing.T) {
	typ := &ditto.Type{}
	err := ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &ditto.Type{
		Type: "test",
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &ditto.Type{
		Value: "test",
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestRegisterGroupAndField(t *testing.T) {
	g := &ditto.Group{
		Name:          "test",
		ValidInfoKeys: nil,
	}
	err := ditto.RegisterGroup(g)
	if err != nil {
		t.Fatal(err)
	}

	err = ditto.RegisterGroup(g)
	if err == nil {
		t.Fatal("error expected")
	}

	typ := &ditto.Type{
		Type:          "test",
		Value:         "empty",
		Group:         g,
		ValidInfoKeys: nil,
	}
	err = ditto.RegisterType(typ)
	if err != nil {
		t.Fatal(err)
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}
