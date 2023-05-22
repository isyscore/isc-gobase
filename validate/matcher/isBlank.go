package matcher

import (
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"strconv"
	"strings"

	"github.com/isyscore/isc-gobase/logger"
)

type IsBlankMatch struct {
	BlackWhiteMatch

	// 是否设置过isNil值
	HaveSet int8

	// 匹配空字符的值
	IsBlank bool
}

func (isBlankMatch *IsBlankMatch) Match(parameterMap map[string]interface{}, _ any, field reflect.StructField, fieldValue any) bool {
	if reflect.TypeOf(fieldValue).Kind() != field.Type.Kind() {
		isBlankMatch.SetBlackMsg("属性 %v 的值不是字符类型", field.Name)
		return false
	}

	if isBlankMatch.IsBlank {
		if fieldValue == "" {
			isBlankMatch.SetBlackMsg("属性 %v 的值为空字符", field.Name)
			return true
		} else {
			isBlankMatch.SetWhiteMsg("属性 %v 的值为非空字符", field.Name)
			return false
		}
	} else {
		if fieldValue != "" {
			isBlankMatch.SetBlackMsg("属性 %v 的值不为空", field.Name)
			return true
		} else {
			isBlankMatch.SetWhiteMsg("属性 %v 的值为空字符", field.Name)
			return false
		}
	}
}

func (isBlankMatch *IsBlankMatch) IsEmpty() bool {
	return isBlankMatch.HaveSet == 0
}

func BuildIsBlankMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errCode, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}

	if !strings.Contains(subCondition, constants.IsBlank) {
		return
	}

	value := "true"
	if strings.Contains(subCondition, constants.EQUAL) {
		index := strings.Index(subCondition, "=")
		value = strings.TrimSpace(subCondition[index+1:])
	}

	if strings.EqualFold("true", value) || strings.EqualFold("false", value) {
		var isBlank bool
		if chgValue, err := strconv.ParseBool(value); err == nil {
			isBlank = chgValue
		} else {
			logger.Error(err.Error())
			return
		}
		addMatcher(objectTypeFullName, objectFieldName, &IsBlankMatch{IsBlank: isBlank, HaveSet: 1}, errCode, errMsg, true)
	}
}
