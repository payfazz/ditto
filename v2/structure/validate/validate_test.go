package validate_test

import (
	"encoding/json"
	"github.com/payfazz/ditto/v2/structure/validate"
	"testing"
)

var validator *validate.Validator

func init() {
	validator = validate.New()
}

func TestValidator_Validate(t *testing.T) {
	type TestCase struct {
		json string
		isError bool
	}

	var cases = []TestCase {
		TestCase{
			json:    formJson,
			isError: false,
		},
		TestCase{
			json:    formJson2,
			isError: false,
		},
		TestCase{
			json:    formJsonFailed,
			isError: true,
		},
		TestCase{
			json:    formJsonFailed2,
			isError: true,
		},
	}

	for i, c := range cases {
		var s map[string]interface{}
		err := json.Unmarshal([]byte(c.json), &s)
		if nil != err {
			t.Log(i+1)
			t.Fatal(err)
		}

		err = validator.Validate(s)
		if (nil == err) == c.isError {
			t.Log(i+1)
			t.Fatal(err)
		}
	}
}

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

var formJson2 = `{
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
			"placeholder": null,
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


var formJsonFailed = `{
	"type": "summary_section_send"
}`

var formJsonFailed2 = `{
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
			"placeholder": null,
			"validations": [
				{
				  "type": "required"
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

