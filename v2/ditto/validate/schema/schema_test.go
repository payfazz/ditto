package schema_test

import (
	"fmt"
	"github.com/payfazz/ditto/v2/ditto/validate/schema"
	"testing"
)

func TestValidate(t *testing.T) {
	result, err := schema.Validate(json1)
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

var json1 = `
{   
    "id": "cnf_ktp_field",
    "type" : "input_text",
    "initial_state": "initial",
    "states" : {
        "initial": {
            "actions" : {   
                "onKeyDown" : [
                    "validateOnKeydown",
                ]
            }
        }, 
        "error_length": {
            "actions" : {   
                "onKeyDown" : [
                    "validateOnKeydown",
                ]
            }
        },
        "error_unique" : {
            "actions" : {   
                "onKeyDown" : [
                    "validateOnKeydown",
                ]
            }
        },
        "success" : {
            "actions" : {   
                "onKeyDown" : [
                    "validateOnKeydown",
                ]
            }
        }
    }
}
`