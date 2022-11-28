package matcher

import (
	"fmt"
	"github.com/isyscore/isc-gobase/constants"
	"reflect"
	"regexp"
	"strings"

	"github.com/isyscore/isc-gobase/logger"
)

type ModelMatch struct {
	BlackWhiteMatch

	isIdCard  bool
	modelName string
	pReg      *regexp.Regexp
}

var modelMap = map[string]*regexp.Regexp{}

func (modelMatch *ModelMatch) Match(_ map[string]interface{}, _ any, field reflect.StructField, fieldValue any) bool {
	if nil == fieldValue {
		return false
	}

	if field.Type.Kind() != reflect.String {
		return false
	}

	// 身份证号单独处理
	if modelMatch.isIdCard {
		if idCardIsValidate(fmt.Sprintf("%v", fieldValue)) {
			modelMatch.SetBlackMsg("属性 %v 的值 %v 符合身份证要求", field.Name, fieldValue)
			return true
		} else {
			modelMatch.SetWhiteMsg("属性 %v 的值 %v 不符合身份证要求", field.Name, fieldValue)
			return false
		}
	} else {
		if modelMatch.pReg.MatchString(fmt.Sprintf("%v", fieldValue)) {
			modelMatch.SetBlackMsg("属性 %v 的值 %v 命中不允许的类型 [%v]", field.Name, fieldValue, modelMatch.modelName)
			return true
		} else {
			modelMatch.SetWhiteMsg("属性 %v 的值 %v 没有命中只允许类型 [%v]", field.Name, fieldValue, modelMatch.modelName)
			return false
		}
	}
}

func (modelMatch *ModelMatch) IsEmpty() bool {
	if modelMatch.isIdCard {
		return false
	}
	return modelMatch.pReg == nil
}

func BuildModelMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}

	if !strings.Contains(subCondition, constants.Model) || !strings.Contains(subCondition, constants.EQUAL) {
		return
	}

	index := strings.Index(subCondition, "=")
	modelKey := strings.TrimSpace(subCondition[index+1:])

	pReg, contain := modelMap[modelKey]
	if !contain && modelKey != constants.IdCard {
		logger.Error("不包含模式%v", modelKey)
		return
	}

	if modelKey == constants.IdCard {
		addMatcher(objectTypeFullName, objectFieldName, &ModelMatch{pReg: pReg, isIdCard: true, modelName: modelKey}, errMsg, true)
	} else {
		addMatcher(objectTypeFullName, objectFieldName, &ModelMatch{pReg: pReg, isIdCard: false, modelName: modelKey}, errMsg, true)
	}
}

func init() {
	// 手机号
	pReg, _ := regexp.Compile("^1(3[0-9]|4[01456879]|5[0-35-9]|6[2567]|7[0-8]|8[0-9]|9[0-35-9])\\d{8}$")
	modelMap[constants.Phone] = pReg

	// 固定电话
	pReg, _ = regexp.Compile("^(([0+]\\d{2,3}-)?(0\\d{2,3})-)(\\d{7,8})(-(\\d{3,}))?$")
	modelMap[constants.FixedPhone] = pReg

	// 邮箱
	pReg, _ = regexp.Compile("^([\\w-_]+(?:\\.[\\w-_]+)*)@[\\w-]+(.[\\w_-]+)+")
	modelMap[constants.MAIL] = pReg

	// IP地址
	pReg, _ = regexp.Compile("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$")
	modelMap[constants.IpAddress] = pReg
}
