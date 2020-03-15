package schema

import (
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
)

func Validate(json string) (*gojsonschema.Result, error) {
	byt, err := ioutil.ReadFile("draft-01.json")
	if nil != err {
		return nil, err
	}

	schemaLoader := gojsonschema.NewStringLoader(string(byt))
	document := gojsonschema.NewStringLoader(json)

	return gojsonschema.Validate(schemaLoader, document)
}
