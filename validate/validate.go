package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/antonmedv/expr"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/validate/matcher"
)

var lock sync.Mutex

type MatchCollector func(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errCode, errMsg string)

type CollectorEntity struct {
	name         string
	infCollector MatchCollector
}

type CheckResult struct {
	Result  bool
	ErrCode string
	ErrMsg  string
}

var checkerEntities []CollectorEntity

/* 核查的标签 */
var matchTagArray = []string{constants.Value, constants.IsBlank, constants.Range, constants.Model, constants.Condition, constants.Regex, constants.Customize}

// Check
// 入参
// 	- object: 检查对象
// 	- fieldNames...: 待检查对象的属性
// 返回值
// 	- bool: 核查结果
// 	- string: 错误code
// 	- string: 错误异常
func Check(object any, fieldNames ...string) (bool, string, string) {
	return CheckWithParameter(map[string]interface{}{}, object, fieldNames...)
}

// Check
// 入参
// 	- parameterMap: 外部参数map，这个一般用于'customize'自定义函数里面的参数用
// 	- object: 检查对象
// 	- fieldNames...: 待检查对象的属性
// 返回值
// 	- bool: 核查结果
// 	- string: 错误code
// 	- string: 错误异常
func CheckWithParameter(parameterMap map[string]interface{}, object interface{}, fieldNames ...string) (bool, string, string) {
	if object == nil {
		return true, "", ""
	}
	objType := reflect.TypeOf(object)
	objValue := reflect.ValueOf(object)

	// 指针类型按照指针类型
	if objType.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		return Check(objValue.Interface(), fieldNames...)
	}

	if objType.Kind() != reflect.Struct {
		return true, "", ""
	}

	// 搜集核查器
	collectCollector(objType)

	ch := make(chan *CheckResult)
	for index, num := 0, objType.NumField(); index < num; index++ {
		field := objType.Field(index)
		fieldValue := objValue.Field(index)

		// 私有字段不处理
		if !isc.IsPublic(field.Name) {
			continue
		}

		// 过滤选择的列
		if !isSelectField(field.Name, fieldNames...) {
			continue
		}

		// 基本类型
		if matcher.IsCheckedKing(fieldValue.Type()) || (fieldValue.Kind() == reflect.Ptr && !fieldValue.Elem().IsValid()) || (fieldValue.Kind() == reflect.Ptr && matcher.IsCheckedKing(fieldValue.Elem().Type())) {
			tagJudge := field.Tag.Get(constants.MATCH)
			if len(tagJudge) == 0 {
				continue
			}

			// 核查结果：任何一个属性失败，则返回失败
			goid.Go(func() {
				check(parameterMap, object, field, fieldValue.Interface(), ch)
			})
			checkResult := <-ch
			if !checkResult.Result {
				close(ch)
				return false, checkResult.ErrCode, checkResult.ErrMsg
			}
		} else if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
			// struct 结构类型
			tagMatch := field.Tag.Get(constants.MATCH)
			if len(tagMatch) == 0 || (len(tagMatch) == 1 && tagMatch != constants.CHECK) {
				continue
			}
			result, errCode, errMsg := Check(fieldValue.Interface())
			if !result {
				return false, errCode, errMsg
			}
		} else if fieldValue.Kind() == reflect.Map || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Map) {
			// map结构
			if fieldValue.Len() == 0 {
				continue
			}

			for mapR := fieldValue.MapRange(); mapR.Next(); {
				mapKey := mapR.Key()
				mapValue := mapR.Value()

				result, errCode, errMsg := Check(mapKey.Interface())
				if !result {
					return false, errCode, errMsg
				}
				result, errCode, errMsg = Check(mapValue.Interface())
				if !result {
					return false, errCode, errMsg
				}
			}
		} else if fieldValue.Kind() == reflect.Array || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Array) {
			// Array 结构
			arrayLen := fieldValue.Len()
			for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
				fieldValueItem := fieldValue.Index(arrayIndex)
				result, errCode, errMsg := Check(fieldValueItem.Interface())
				if !result {
					return false, errCode, errMsg
				}
			}
		} else if fieldValue.Kind() == reflect.Slice || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Slice) {
			// Slice 结构
			tagJudge := field.Tag.Get(constants.MATCH)
			if len(tagJudge) == 0 {
				continue
			}

			// 核查结果：任何一个属性失败，则返回失败
			goid.Go(func() {
				check(parameterMap, object, field, fieldValue.Interface(), ch)
			})
			checkResult := <-ch
			if !checkResult.Result {
				close(ch)
				return false, checkResult.ErrCode, checkResult.ErrMsg
			}

			arrayLen := fieldValue.Len()
			for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
				fieldValueItem := fieldValue.Index(arrayIndex)
				result, errCode, errMsg := Check(fieldValueItem.Interface())
				if !result {
					return false, errCode, errMsg
				}
			}
		}
	}
	close(ch)
	return true, "", ""
}

// 搜集核查器
func collectCollector(objType reflect.Type) {
	objectFullName := objType.String()

	/* 搜集过则不再搜集 */
	if _, contain := matcher.MatchMap[objectFullName]; contain {
		return
	}

	lock.Lock()
	/* 搜集过则不再搜集 */
	if _, contain := matcher.MatchMap[objectFullName]; contain {
		return
	}

	doCollectCollector(objType)
	lock.Unlock()
}

func doCollectCollector(objType reflect.Type) {
	// 基本类型不需要搜集
	if matcher.IsCheckedKing(objType) {
		return
	}

	// 指针类型按照指针类型
	if objType.Kind() == reflect.Ptr {
		doCollectCollector(objType.Elem())
		return
	}

	if objType.Kind() != reflect.Struct {
		return
	}

	objectFullName := objType.String()
	for fieldIndex, num := 0, objType.NumField(); fieldIndex < num; fieldIndex++ {
		field := objType.Field(fieldIndex)
		fieldKind := field.Type.Kind()

		// 不可访问字段不处理
		if !isc.IsPublic(field.Name) {
			continue
		}

		if fieldKind == reflect.Ptr {
			fieldKind = field.Type.Elem().Kind()
		}

		// 禁用
		tagMatch := field.Tag.Get(constants.Disable)
		if len(tagMatch) != 0 && tagMatch == "true" {
			continue
		}

		// 基本类型
		if matcher.IsCheckedKing(field.Type) {
			// 错误码信息
			errMsg := field.Tag.Get(constants.ErrMsg)

			// 错误码code
			errCode := field.Tag.Get(constants.ErrCode)

			// match
			tagMatch := field.Tag.Get(constants.MATCH)
			if len(tagMatch) == 0 {
				continue
			}

			if _, contain := matcher.MatchMap[objectFullName][field.Name]; !contain {
				addMatcher(objectFullName, fieldKind, field.Name, tagMatch, errCode, errMsg)
			}

			// accept
			tagAccept := field.Tag.Get(constants.Accept)
			if len(tagMatch) == 0 {
				continue
			}

			if _, contain := matcher.MatchMap[objectFullName][field.Name]; contain {
				addCollector(objectFullName, fieldKind, field.Name, constants.Accept, tagAccept, errCode, errMsg)
			}
		} else if fieldKind == reflect.Struct {
			// struct 结构类型
			tagMatch := field.Tag.Get(constants.MATCH)
			if len(tagMatch) == 0 || (len(tagMatch) == 1 && tagMatch != constants.CHECK) {
				continue
			}

			doCollectCollector(field.Type)
		} else if fieldKind == reflect.Map {
			// Map 结构
			doCollectCollector(field.Type.Key())
			doCollectCollector(field.Type.Elem())
		} else if fieldKind == reflect.Array {
			// Array 结构
			doCollectCollector(field.Type.Elem())
		} else if fieldKind == reflect.Slice {
			// Slice 结构

			// 错误码信息
			errMsg := field.Tag.Get(constants.ErrMsg)

			// 错误码code
			errCode := field.Tag.Get(constants.ErrCode)

			// match
			tagMatch := field.Tag.Get(constants.MATCH)
			if len(tagMatch) == 0 {
				continue
			}

			if _, contain := matcher.MatchMap[objectFullName][field.Name]; !contain {
				addMatcher(objectFullName, fieldKind, field.Name, tagMatch, errCode, errMsg)
			}

			// accept
			tagAccept := field.Tag.Get(constants.Accept)
			if len(tagMatch) == 0 {
				continue
			}

			if _, contain := matcher.MatchMap[objectFullName][field.Name]; !contain {
				addCollector(objectFullName, fieldKind, field.Name, constants.Accept, tagAccept, errCode, errMsg)
			}

			doCollectCollector(field.Type.Elem())
		} else {
			// Uintptr 类型不处理
		}
	}
}

// 是否是选择的列，没有选择也认为是选择的
func isSelectField(fieldName string, fieldNames ...string) bool {
	if len(fieldNames) == 0 {
		return true
	}
	for _, name := range fieldNames {
		// 不区分大小写
		if strings.EqualFold(name, fieldName) {
			return true
		}
	}
	return false
}

// 搜集处理器，对于有一些空格的也进行单独处理
func addMatcher(objectFullName string, fieldKind reflect.Kind, fieldName string, matchJudge string, errCode, errMsg string) {
	var subStrIndexes []int
	for _, tag := range matchTagArray {
		index := strings.Index(matchJudge, tag)
		if index != -1 {
			subStrIndexes = append(subStrIndexes, index)
		}
	}
	sort.Ints(subStrIndexes)

	lastIndex := 0
	for _, subIndex := range subStrIndexes {
		if lastIndex == subIndex {
			continue
		}
		subJudgeStr := matchJudge[lastIndex:subIndex]
		buildChecker(objectFullName, fieldKind, fieldName, constants.MATCH, subJudgeStr, errCode, errMsg)
		lastIndex = subIndex
	}

	subJudgeStr := matchJudge[lastIndex:]
	buildChecker(objectFullName, fieldKind, fieldName, constants.MATCH, subJudgeStr, errCode, errMsg)
}

// 添加搜集器
func addCollector(objectFullName string, fieldKind reflect.Kind, fieldName string, tagName string, matchJudge string, errCode, errMsg string) {
	buildChecker(objectFullName, fieldKind, fieldName, tagName, matchJudge, errCode, errMsg)
}

func buildChecker(objectFullName string, fieldKind reflect.Kind, fieldName string, tagName string, subStr string, errCode, errMsg string) {
	for _, entity := range checkerEntities {
		entity.infCollector(objectFullName, fieldKind, fieldName, tagName, subStr, errCode, errMsg)
	}
}

func check(parameterMap map[string]interface{}, object any, field reflect.StructField, fieldRelValue any, ch chan *CheckResult) {
	objectType := reflect.TypeOf(object)

	if fieldMatcher, contain := matcher.MatchMap[objectType.String()][field.Name]; contain {
		accept := fieldMatcher.Accept
		errMsgOfMatcherProgram := fieldMatcher.ErrMsgProgram
		errCodeOfMatcher := fieldMatcher.ErrCode
		matchers := fieldMatcher.Matchers

		// 黑名单，而且匹配到，则核查失败
		if !accept {
			if matchResult, _errCode, _errMsg := judgeMatch(matchers, parameterMap, object, field, fieldRelValue, accept); matchResult {
				errMsgFinal := ""
				errCodeFinal := errCodeOfMatcher
				if errMsgOfMatcherProgram != nil {
					env := map[string]any{
						"sprintf": fmt.Sprintf,
						"root":    object,
						"current": fieldRelValue,
					}

					output, err := expr.Run(errMsgOfMatcherProgram, env)
					if err != nil {
						logger.Error(err.Error())
						ch <- &CheckResult{Result: false, ErrMsg: err.Error()}
						return
					}

					result := fmt.Sprintf("%v", output)

					errMsgFinal = result
				} else {
					errMsgFinal = _errMsg
				}

				if _errCode != "" {
					errCodeFinal = _errCode
				}

				ch <- &CheckResult{Result: false, ErrCode: errCodeFinal, ErrMsg: errMsgFinal}
				return
			}
		}

		// 白名单，没有匹配到，则核查失败
		if accept {
			if matchResult, _errCode, _errMsg := judgeMatch(matchers, parameterMap, object, field, fieldRelValue, accept); !matchResult {
				errMsgFinal := ""
				errCodeFinal := errCodeOfMatcher
				if errMsgOfMatcherProgram != nil {
					env := map[string]any{
						"sprintf": fmt.Sprintf,
						"root":    object,
						"current": fieldRelValue,
					}

					output, err := expr.Run(errMsgOfMatcherProgram, env)
					if err != nil {
						logger.Error(err.Error())
						ch <- &CheckResult{Result: false, ErrMsg: err.Error()}
						return
					}

					result := fmt.Sprintf("%v", output)
					errMsgFinal = result
				} else {
					errMsgFinal = _errMsg
				}

				if _errCode != "" {
					errCodeFinal = _errCode
				}
				ch <- &CheckResult{Result: false, ErrCode: errCodeFinal, ErrMsg: errMsgFinal}
				return
			}
		}
	}
	ch <- &CheckResult{Result: true}
	return
}

// 任何一个匹配上，则返回true，都没有匹配上则返回false
func judgeMatch(matchers []*matcher.Matcher, parameterMap map[string]interface{}, object any, field reflect.StructField, fieldValue any, accept bool) (bool, string, string) {
	var errMsgArray []string
	var errCode string
	for _, match := range matchers {
		if (*match).IsEmpty() {
			continue
		}

		matchResult := (*match).Match(parameterMap, object, field, fieldValue)
		if matchResult {
			if !accept {
				errMsgArray = append(errMsgArray, (*match).GetBlackMsg())
				errCode = (*match).GetErrCode()
			} else {
				errMsgArray = []string{}
			}
			return true, errCode, arraysToString(errMsgArray)
		} else {
			if accept {
				errMsgArray = append(errMsgArray, (*match).GetWhitMsg())
				errCode = (*match).GetErrCode()
			}
		}
	}
	return false, errCode, arraysToString(errMsgArray)
}

func RegisterCustomize(funName string, fun any) {
	matcher.RegisterCustomize(funName, fun)
}

// 包的初始回调
func init() {
	/* 匹配后是否接受 */
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Accept, matcher.CollectAccept})

	/* 搜集匹配器 */
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Value, matcher.BuildValuesMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.IsBlank, matcher.BuildIsBlankMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.IsUnBlank, matcher.BuildIsUnBlankMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Range, matcher.BuildRangeMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Model, matcher.BuildModelMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Condition, matcher.BuildConditionMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Customize, matcher.BuildCustomizeMatcher})
	checkerEntities = append(checkerEntities, CollectorEntity{constants.Regex, matcher.BuildRegexMatcher})
}

func arraysToString(dataArray []string) string {
	if len(dataArray) == 1 {
		return dataArray[0]
	}
	myValue, _ := json.Marshal(dataArray)
	return string(myValue)
}
