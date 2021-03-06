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

type FieldRule func(value interface{}, rule interface{}) bool

var rules = make(map[string]FieldRule)

func GetRule(key string) FieldRule {
	val, ok := rules[key]
	if !ok {
		return nil
	}
	return val
}

func RegisterRule(tag string, fn FieldRule) error {
	if len(tag) == 0 {
		return ErrKeyEmpty
	}

	if fn == nil {
		return ErrFuncEmpty
	}

	rules[tag] = fn
	return nil
}

func init() {
	_ = RegisterRule("required", func(value interface{}, rule interface{}) bool {
		if value == nil {
			return false
		}

		if value.(string) == "" {
			return false
		}

		return true
	})

	_ = RegisterRule("text_length_between", func(value interface{}, rule interface{}) bool {
		splitBetween := strings.Split(rule.(string), ",")
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

	_ = RegisterRule("age_between", func(value interface{}, rule interface{}) bool {
		splitBetween := strings.Split(rule.(string), ",")
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

	_ = RegisterRule("date_between", func(value interface{}, rule interface{}) bool {
		splitBetween := strings.Split(rule.(string), ",")

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

	_ = RegisterRule("regex", func(value interface{}, rule interface{}) bool {
		valueString, ok := value.(string)
		if !ok {
			return false
		}

		validationVal, _ := url.QueryUnescape(rule.(string))
		match, _ := regexp.MatchString(validationVal, valueString)
		if !match {
			return false
		}

		return true
	})
}
