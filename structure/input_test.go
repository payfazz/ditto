package structure_test

import (
	"encoding/json"
	"github.com/payfazz/ditto/structure"
	"testing"
)

func TestInput(t *testing.T) {
	var s map[string]interface{}
	err := json.Unmarshal([]byte(formJson), &s)
	if nil != err {
		t.Fatal(err)
	}

	root, err := structure.NewSectionFromMap(s)
	if nil != err {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(inputJson), &data)
	if nil != err {
		t.Fatal(err)
	}

	values := data["values"].(map[string]interface{})

	structure.InputForm(root, values)
	t.Log(values)
}

var inputJson = `{
  "values": {
    "name": {
      "value": "",
      "type": "text"
    }
  }
}`

var formJson = `{
	"id": "task_form_section",
	"title": "Nama",
	"type": "summary_section_send",
	"child": [
		{
			"id": "name",
			"info": {
				"icon": "https://content.payfazz.com/object/e1a14d808797311aaba5f164d50e30fb7ec037a15b07f6e687f31eaa49b4c60f"
			},
			"type": "text",
			"title": "Nama",
			"description": "Nama",
			"placeholder": "Nama",
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