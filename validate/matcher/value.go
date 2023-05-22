package matcher

import (
	"fmt"
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"strings"

	"github.com/isyscore/isc-gobase/logger"
)

type ValueMatch struct {
	BlackWhiteMatch
	Values []any
}

func (valueMatch *ValueMatch) Match(_ map[string]interface{}, _ any, field reflect.StructField, fieldValue any) bool {
	values := valueMatch.Values

	fieldRelValue := fieldValue
	if reflect.TypeOf(fieldValue).Kind() == reflect.Ptr {
		fieldRelValue = reflect.ValueOf(fieldValue).Elem()
	}

	for _, value := range values {
		if fmt.Sprintf("%v", value) == fmt.Sprintf("%v", fieldRelValue) {
			valueMatch.SetBlackMsg("属性 %v 的值 %v 位于禁用值 %v 中", field.Name, fieldRelValue, values)
			return true
		}
	}
	valueMatch.SetWhiteMsg("属性 %v 的值 %v 不在只可用列表 %v 中", field.Name, fieldRelValue, values)
	return false
}

func (valueMatch *ValueMatch) IsEmpty() bool {
	return len(valueMatch.Values) == 0
}

func BuildValuesMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errCode, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}
	if !strings.Contains(subCondition, constants.Value) || !strings.Contains(subCondition, constants.EQUAL) {
		return
	}

	index := strings.Index(subCondition, "=")
	value := subCondition[index+1:]

	var availableValues []any
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		value = value[1 : len(value)-1]
		for _, subValue := range strings.Split(value, ",") {
			subValue = strings.TrimSpace(subValue)
			if chgValue, err := Cast(fieldKind, subValue); err == nil {
				availableValues = append(availableValues, chgValue)
			} else {
				logger.Error(err.Error())
				continue
			}
		}
	} else {
		value = strings.TrimSpace(value)
		if chgValue, err := Cast(fieldKind, value); err == nil {
			availableValues = append(availableValues, chgValue)
		} else {
			logger.Error(err.Error())
			return
		}
	}
	addMatcher(objectTypeFullName, objectFieldName, &ValueMatch{Values: availableValues}, errCode, errMsg, true)
}
