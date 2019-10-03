package xml2json_test

import (
	"encoding/json"
	"github.com/payfazz/ditto/xml2json"
	"reflect"
	"testing"
)

func TestXML2JSON(t *testing.T) {
	jsonResult, err := xml2json.XMLToDittoJSON(xmlString)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(jsonResult)

	isEquals, err := JSONBytesEqual([]byte(expectedJSON), []byte(jsonResult))
	if err != nil {
		t.Fatal(err)
	}

	if !isEquals {
		t.Fatal("true expected")
	}
}

func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

var xmlString = `
<SummarySectionSend id="test" description="test" title="test">
<Text id="a" title="a title" description="a desc" validate="required#Harus diisi|regex:%5E%5B%5Cw%20%5D%7B2%2C100%7D%24#Pastikan nama toko alphanumeric"/>
</SummarySectionSend>
`

var expectedJSON = `{
	"id": "test",
	"type":"summary_section_send",
	"description": "test",
	"title":"test",
	"child": [
		{
			"id": "a",
			"type": "text",
			"title": "a title",
			"description": "a desc",
			"validations": [
				{
				  "type": "required",
				  "error_message": "Harus diisi"
				},
				{
				  "type": "regex",
				  "value": "%5E%5B%5Cw%20%5D%7B2%2C100%7D%24",
				  "error_message": "Pastikan nama toko alphanumeric"
				}
		  	]
		}
	]
}`
