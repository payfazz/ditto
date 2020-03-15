package schema_test

import (
	"fmt"
	"github.com/payfazz/ditto/v2/ditto/validate/schema"
	"io/ioutil"
	"testing"
)

func TestValidate(t *testing.T) {
	byt, err := ioutil.ReadFile("example.json")
	result, err := schema.Validate(string(byt))
	if nil != err {
		t.Fatal(err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		t.Fatal("valid expected")
	}
}

var json1 = ``