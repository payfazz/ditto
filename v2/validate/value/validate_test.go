package value_test

import (
	"encoding/json"
	"github.com/payfazz/ditto/v2/validate"
	"testing"
)

var validator *validate.Validator

func init() {
	validator = validate.New()
}

func TestValidator_ExtractField(t *testing.T) {
	var s map[string]interface{}
	err := json.Unmarshal([]byte(formJson2), &s)
	if nil != err {
		t.Fatal(err)
	}

	m, err := validator.ExtractField(s)
	if nil != err {
		t.Fatal(err)
	}

	t.Log(m)
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
		},
		{
			"id": "b",
			"type": "text",
			"title": "b title",
			"description": "b desc",
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

