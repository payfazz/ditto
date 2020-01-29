package value

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrKeyEmpty = errors.New("function key cannot be empty")
var ErrFuncEmpty = errors.New("function cannot be empty")

type FieldValidator func(value interface{}, rule string) bool

var validators = make(map[string]FieldValidator)

func GetValidator(key string) FieldValidator {
	val, ok := validators[key]
	if !ok {
		return nil
	}
	return val
}

func RegisterValidator(tag string, fn FieldValidator) error {
	if len(tag) == 0 {
		return ErrKeyEmpty
	}

	if fn == nil {
		return ErrFuncEmpty
	}

	validators[tag] = fn
	return nil
}

func init() {
	_ = RegisterValidator("required", func(value interface{}, rule string) bool {
		if value == nil {
			return false
		}
		valueObj, ok := value.(map[string]interface{})
		if !ok {
			return false
		}
		if valueObj["value"] == nil {
			return false
		}

		if valueObj["value"].(string) == "" {
			return false
		}

		return true
	})

	_ = RegisterValidator("text_length_between", func(value interface{}, rule string) bool {
		splitBetween := strings.Split(rule, ",")
		min, _ := strconv.Atoi(splitBetween[0])
		max, _ := strconv.Atoi(splitBetween[1])

		valueObj := value.(map[string]interface{})
		valueString, ok := valueObj["value"].(string)
		if !ok {
			return false
		}

		if len(valueString) < min || len(valueString) > max {
			return false
		}

		return true
	})

	_ = RegisterValidator("age_between", func(value interface{}, rule string) bool {
		splitBetween := strings.Split(rule, ",")
		min, _ := strconv.Atoi(splitBetween[0])
		max, _ := strconv.Atoi(splitBetween[1])

		valueObj := value.(map[string]interface{})
		valueString, ok := valueObj["value"].(string)
		if !ok {
			return false
		}

		birthDate, _ := time.Parse("02-01-2006", valueString)
		now := time.Now()
		years := now.Year() - birthDate.Year()
		if now.YearDay() < birthDate.YearDay() {
			years--
		}
		if years < min || years > max {
			return false
		}

		return true
	})

	_ = RegisterValidator("date_between", func(value interface{}, rule string) bool {
		splitBetween := strings.Split(rule, ",")

		valueObj := value.(map[string]interface{})
		valueString, ok := valueObj["value"].(string)
		if !ok {
			return false
		}

		date, _ := time.Parse("02-01-2006", valueString)
		min, err := time.Parse("02-01-2006", splitBetween[0])
		if nil != err {
			return false
		}

		max, err := time.Parse("02-01-2006", splitBetween[1])
		if nil != err {
			return false
		}

		if date.UTC().Unix() < min.UTC().Unix() || date.UTC().Unix() > max.UTC().Unix() {
			return false
		}

		return true
	})

	_ = RegisterValidator("regex", func(value interface{}, rule string) bool {
		valueObj := value.(map[string]interface{})
		valueString, ok := valueObj["value"].(string)
		if !ok {
			return false
		}

		validationVal, _ := url.QueryUnescape(rule)
		match, _ := regexp.MatchString(validationVal, valueString)
		if !match {
			return false
		}

		return true
	})
}
