package matcher

import (
	"fmt"
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"regexp"
	"strings"
)

type RegexMatch struct {
	BlackWhiteMatch

	Reg *regexp.Regexp
}

func (regexMatch *RegexMatch) Match(_ any, field reflect.StructField, fieldValue any) bool {
	if regexMatch.Reg.MatchString(fmt.Sprintf("%v", fieldValue)) {
		regexMatch.SetBlackMsg("属性 %v 的值 %v 命中禁用的正则表达式 %v ", field.Name, fieldValue, regexMatch.Reg.String())
		return true
	} else {
		regexMatch.SetWhiteMsg("属性 %v 的值 %v 没命中只允许的正则表达式 %v ", field.Name, fieldValue, regexMatch.Reg.String())
	}
	return false
}

func (regexMatch *RegexMatch) IsEmpty() bool {
	return regexMatch.Reg == nil
}

func BuildRegexMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}
	if !strings.Contains(subCondition, constants.Regex) || !strings.Contains(subCondition, constants.EQUAL) {
		return
	}

	index := strings.Index(subCondition, "=")
	value := subCondition[index+1:]

	reg, err := regexp.Compile(value)
	if err != nil {
		return
	}
	addMatcher(objectTypeFullName, objectFieldName, &RegexMatch{Reg: reg}, errMsg, true)
}
