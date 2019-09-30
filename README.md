# ditto

#### Table of Contents  
1. [Validation](#validation)
2. [Group](#group)
3. [Type](#type)
4. [Create Form](#create-form)
5. [Validate Input](#validate-input)

<a name="validation"/>

## Validation

Default validations:
- required
- text_length_between
- age_between
- date_between
- regex

Register new validation

```go
import "ditto"

ditto.RegisterValidator("file_extension", func(value interface{}, fieldVal ditto.FieldValidation) bool {
    validExtensions := strings.Split(fieldVal.Value, "|")

    valueObj := value.(map[string]interface{})
    valueString, ok := valueObj["value"].(string)
    if !ok {
        return false
    }

    valueSplit := strings.Split(valueString, ".")
    valueLen := len(valueSplit)
    if valueLen <= 1 {
        return false
    }
    
    ext := valueSplit(valueLen - 1)

    for _, validExt := range validExtensions {
        if ext == validExt {
            return true
        }
    }

    return false
})
```

<a name="group"/>

## Group

Default groups:
- section
- section_field
- text
- file
- list


<a name="type"/>

## Type

Default types:
- summary_field
- nextable_section
- nextable_form
- nextable_field
- summary_section_send
- summary_section_save
- text_multiline
- text_numeric
- photo_camera
- list
- text
- date
- searchable_list
- object_searchable_list
- normal_list

<a name="create-form"/>

## Create Form

Use **NewSectionFromMap** function to generate form. Generated form is started from root **Section**.

```go
import "ditto"

var formMap map[string]interface{}
err := json.Unmarshal([]byte(jsonData), &formMap)
if nil != err {
    t.Fatal(err)
}

root, err := ditto.NewSectionFromMap(s)
if nil != err {
    t.Fatal(err)
}
```

<a name="validate-input"/>

## Validate Input

Use **ValidateFormInput** function to validate input.

```go
var formMap map[string]interface{}
_ := json.Unmarshal([]byte(jsonData), &formMap)
root, _ := ditto.NewSectionFromMap(s)

var inputMap map[string]interface{}
_ = json.Unmarshal([]byte(inputJson), &inputMap)

values := data["values"].(map[string]interface{})
result, err := ditto.ValidateFormInput(root, values)
```
