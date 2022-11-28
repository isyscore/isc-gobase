package matcher

import (
	"fmt"
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/isyscore/isc-gobase/logger"
)

type ConditionMatch struct {
	BlackWhiteMatch

	expression string
	Program    *vm.Program
}

func (conditionMatch *ConditionMatch) Match(_ map[string]interface{}, object any, field reflect.StructField, fieldValue any) bool {
	env := map[string]any{
		"root":    object,
		"current": fieldValue,
	}

	output, err := expr.Run(conditionMatch.Program, env)
	if err != nil {
		logger.Error("表达式 %v 执行失败: %v", conditionMatch.expression, err.Error())
		return false
	}

	result, err := CastBool(fmt.Sprintf("%v", output))
	if err != nil {
		return false
	}

	if result {
		conditionMatch.SetBlackMsg("属性 %v 的值 %v 命中禁用条件 [%v] ", field.Name, fieldValue, conditionMatch.expression)
	} else {
		conditionMatch.SetWhiteMsg("属性 %v 的值 %v 不符合条件 [%v] ", field.Name, fieldValue, conditionMatch.expression)
	}
	return result
}

func (conditionMatch *ConditionMatch) IsEmpty() bool {
	return conditionMatch.Program == nil
}

func BuildConditionMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}
	if !strings.Contains(subCondition, constants.Condition) || !strings.Contains(subCondition, constants.EQUAL) {
		return
	}

	index := strings.Index(subCondition, constants.EQUAL)
	expression := subCondition[index+1:]

	if expression == "" {
		return
	}

	tree, err := parser.Parse(rmvWell(expression))
	if err != nil {
		logger.Error("脚本：%v 解析异常：%v", expression, err.Error())
		return
	}

	program, err := compiler.Compile(tree, nil)
	if err != nil {
		logger.Error("脚本: %v 编译异常：%v", expression, err.Error())
		return
	}
	addMatcher(objectTypeFullName, objectFieldName, &ConditionMatch{Program: program, expression: expression}, errMsg, true)
}
