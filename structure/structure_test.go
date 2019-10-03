package structure_test

import (
	"fmt"
	"github.com/payfazz/ditto/structure"
	"testing"
)

func TestField(t *testing.T) {
	c1, err := structure.CreateComponent("text")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c1.RequiredKeys())

	c2, err := structure.CreateComponent("file")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c2.RequiredKeys())

	c3, err := structure.CreateComponent("list")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c3.RequiredKeys())


	_, err = structure.CreateComponent("a")
	if nil == err {
		t.Fatal("error expected")
	}
}
