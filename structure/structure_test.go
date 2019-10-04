package structure_test

import (
	"fmt"
	"github.com/payfazz/ditto/structure"
	"github.com/payfazz/ditto/structure/field"
	"testing"
)

func TestRequiredAttrsField(t *testing.T) {
	c1, err := structure.CreateComponent("text")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c1.RequiredAttrs())

	c2, err := structure.CreateComponent("file")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c2.RequiredAttrs())

	c3, err := structure.CreateComponent("list")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(c3.RequiredAttrs())

	_, err = structure.CreateComponent("a")
	if nil == err {
		t.Fatal("error expected")
	}
}

func TestExtractValidationField(t *testing.T) {
	c1, err := structure.CreateComponent("text")
	if nil != err {
		t.Fatal(err)
	}

	attrs := map[string]interface{}{
		"id":          "test",
		"description": "test",
		"title":       "test",
		"type":        "test",
		"validations": []interface{}{
			map[string]interface{}{
				"type": "required",
				"error_message": "",
			},
		},
	}

	err = c1.FillStruct(attrs)
	if nil != err {
		t.Fatal(err)
	}

	t.Logf("%+v", (c1.(*field.Text)).Validations)
}
