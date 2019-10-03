package xml2json

import (
	"encoding/json"
	xj "github.com/basgys/goxml2json"
	"github.com/iancoleman/strcase"
	"strings"
)

func XMLToDittoJSON(xml string) (string, error) {
	reader := strings.NewReader(xml)

	jsonResult, err := xj.Convert(reader)
	if nil != err {
		return "", err
	}

	var s map[string]interface{}
	err = json.Unmarshal(jsonResult.Bytes(), &s)
	if nil != err {
		return "", err
	}

	result := parseJSONtoDitto(s)

	m, err := json.Marshal(result)
	if nil != err {
		return "", nil
	}

	return string(m), nil
}

func parseJSONtoDitto(s map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for k, v := range s {
		result["type"] = strcase.ToSnake(k)

		member := v.(map[string]interface{})
		childs := make([]map[string]interface{}, 0)
		for k2, v2 := range member {
			if k2[0] == '-' {
				if k2[1:] == "validate" {
					result["validations"] = parseValidate(v2)
					continue
				}

				result[k2[1:]] = v2
				continue
			}

			child := parseJSONtoDitto(map[string]interface{}{
				k2: v2,
			})
			childs = append(childs, child)
		}

		if len(childs) > 0 {
			result["child"] = childs
		}
	}

	return result
}

func parseValidate(v interface{}) []map[string]interface{} {
	validations := make([]map[string]interface{}, 0)
	s := v.(string)
	arrStr := strings.Split(s, "|")
	for _, str := range arrStr {
		result := make(map[string]interface{})
		splits := strings.Split(str, "#")
		result["error_message"] = ""
		if len(splits) >= 1 {
			result["error_message"] = splits[1]
		}

		typeAndValue := strings.Split(splits[0], ":")

		result["type"] = typeAndValue[0]
		if len(typeAndValue) >= 2 {
			result["value"] = typeAndValue[1]
		}

		validations = append(validations, result)
	}

	return validations
}
