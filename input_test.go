package ditto_test

import (
	"encoding/json"
	"github.com/payfazz/ditto"
	"testing"
)

func TestInput(t *testing.T) {
	var s map[string]interface{}
	err := json.Unmarshal([]byte(formJson), &s)
	if nil != err {
		t.Fatal(err)
	}

	root, err := ditto.NewSectionFromMap(s)
	if nil != err {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(inputJson), &data)
	if nil != err {
		t.Fatal(err)
	}

	values := data["values"].(map[string]interface{})
	t.Log(values)

	result, err := ditto.ValidateFormInput(root, values)

	m, err := json.Marshal(result)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(string(m))
	t.Log(result.Structure)
}

var inputJson = `{
  "values": {
    "name": {
      "value": "@@@",
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
