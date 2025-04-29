package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	validateStrLenErr  = "должно содержать %d символов"
	validateStrRegxErr = "не валидно"
	validateInListErr  = "должно быть одним из: %s"
	validateIntMinErr  = "должно быть меньше %d"
	validateIntMaxErr  = "должно быть больше %d"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type validateRule struct {
	rule   string
	params string
}

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}

	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Поле %s %s\n", err.Field, err.Err.Error()))
	}

	return sb.String()
}

func Validate(v interface{}) error {
	var vErrors ValidationErrors
	var rules []validateRule

	t := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)

	for i := 0; i < vo.NumField(); i++ {
		rules = parseRulesFromTag(t.Field(i).Tag.Get("validate"))
		validateField(&vErrors, t.Field(i), vo.Field(i), rules)
	}

	if len(vErrors) > 0 {
		return errors.New(vErrors.Error())
	}

	return nil
}

// Распарсить правило из тега.
func parseRulesFromTag(validateTag string) []validateRule {
	var rules, ruleInfo []string
	var result []validateRule

	rules = strings.Split(validateTag, "|")

	for _, ruleItem := range rules {
		ruleInfo = strings.SplitN(ruleItem, ":", 2)

		if len(ruleInfo) == 2 {
			result = append(result, validateRule{rule: ruleInfo[0], params: ruleInfo[1]})
		}
	}

	return result
}

// Валидация поля.
func validateField(vErrors *ValidationErrors, fData reflect.StructField, fValue reflect.Value, rules []validateRule) {
	for _, rule := range rules {
		validateByType(vErrors, fValue.Kind().String(), fData.Name, fValue, rule)
	}
}

// Валидация по типу поля.
func validateByType(vErrors *ValidationErrors, fType, fName string, fValue reflect.Value, rule validateRule) {
	var err error

	switch fType {
	case "string":
		err = validateString(fValue.String(), rule)
	case "int":
		err = validateInteger(int(fValue.Int()), rule)
	case "slice":
		switch fValue.Type().String() {
		case "[]string":
			for _, v := range fValue.Interface().([]string) {
				validateByType(vErrors, "string", fName, reflect.ValueOf(v), rule)
			}
		case "[]int":
			for _, v := range fValue.Interface().([]int) {
				validateByType(vErrors, "int", fName, reflect.ValueOf(v), rule)
			}
		}
	}

	if err != nil {
		addError(vErrors, fName, err)
	}
}

// Добавить ошибку в массив.
func addError(vErrors *ValidationErrors, fName string, err error) {
	*vErrors = append(*vErrors, ValidationError{
		Field: fName,
		Err:   err,
	})
}

// Валидация строк.
func validateString(value string, rule validateRule) error {
	switch rule.rule {
	case "len":
		ruleLen, err := strconv.Atoi(rule.params)
		if err != nil {
			return err
		}

		if len(value) != ruleLen {
			return fmt.Errorf(validateStrLenErr, ruleLen)
		}
	case "regexp":
		matchString, err := regexp.MatchString(rule.params, value)
		if err != nil {
			return err
		}

		if !matchString {
			return errors.New(validateStrRegxErr)
		}
	case "in":
		list := strings.Split(rule.params, ",")

		if !slices.Contains(list, value) {
			return fmt.Errorf(validateInListErr, rule.params)
		}
	}

	return nil
}

// Валидация чисел.
func validateInteger(value int, rule validateRule) error {
	switch rule.rule {
	case "min":
		minVal, err := strconv.Atoi(rule.params)
		if err != nil {
			return err
		}

		if value < minVal {
			return fmt.Errorf(validateIntMinErr, minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(rule.params)
		if err != nil {
			return err
		}

		if value > maxVal {
			return fmt.Errorf(validateIntMaxErr, maxVal)
		}

	case "in":
		list := strings.Split(rule.params, ",")
		var listInt []int
		for _, v := range list {
			listItem, err := strconv.Atoi(v)
			if err != nil {
				return err
			}

			listInt = append(listInt, listItem)
		}

		if !slices.Contains(listInt, value) {
			return fmt.Errorf(validateInListErr, rule.params)
		}
	}

	return nil
}
