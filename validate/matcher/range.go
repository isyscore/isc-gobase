package matcher

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/isyscore/isc-gobase/constant"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/validate/util"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RangeMatch struct {
	BlackWhiteMatch

	RangeExpress string
	Script       string
	Begin        interface{}
	End          interface{}
	BeginNow     bool
	EndNow       bool
	Program      *vm.Program
}

type RangeEntity struct {
	beginAli    string
	begin       interface{}
	end         interface{}
	endAli      string
	dateFlag    bool
	dynamicTime bool
	beginNow    bool
	endNow      bool
}

type DynamicTimeNum struct {
	plusOrMinus bool
	years       int
	months      int
	days        int
	hours       int
	minutes     int
	seconds     int
}

type Predicate func(subCondition string) bool

// []：时间或者数字范围匹配
var rangeRegex = regexp.MustCompile("^(\\(|\\[)(.*)(,|，)(\\s)*(.*)(\\)|\\])$")

// digitRegex 全是数字匹配（整数，浮点数，0，负数）
var digitRegex = regexp.MustCompile("^(0)|^[-+]?([1-9]+\\d*|0\\.(\\d*)|[1-9]\\d*\\.(\\d*))$")

// 时间的前后计算匹配：(-|+)yMd(h|H)msS
var timePlusRegex = regexp.MustCompile("^([-+])?(\\d*y)?(\\d*M)?(\\d*d)?(\\d*H|\\d*h)?(\\d*m)?(\\d*s)?$")

func (rangeMatch *RangeMatch) Match(object interface{}, field reflect.StructField, fieldValue interface{}) bool {
	env := map[string]interface{}{
		"begin": rangeMatch.Begin,
		"end":   rangeMatch.End,
	}

	fieldKind := field.Type.Kind()
	if util.IsCheckNumber(fieldKind) {
		env["value"] = fieldValue
	} else if fieldKind == reflect.String {
		env["value"] = len(fmt.Sprintf("%v", fieldValue))
	} else if fieldKind == reflect.Slice {
		env["value"] = reflect.ValueOf(fieldValue).Len()
	} else if field.Type.String() == "time.Time" {
		env["value"] = fieldValue.(time.Time).UnixNano()
		if rangeMatch.BeginNow {
			env["begin"] = time.Now().UnixNano()
		} else if rangeMatch.EndNow {
			env["end"] = time.Now().UnixNano()
		}
	} else {
		return true
	}

	output, err := expr.Run(rangeMatch.Program, env)
	if err != nil {
		logger.Error("脚本 %v 执行失败: %v", rangeMatch.Script, err.Error())
		return false
	}

	result, err := util.CastBool(fmt.Sprintf("%v", output))
	if err != nil {
		return false
	}

	if result {
		if field.Type.Kind() == reflect.String {
			if len(fmt.Sprintf("%v", fieldValue)) > 1024 {
				rangeMatch.SetBlackMsg("属性 %v 值的字符串长度位于禁用的范围 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
			} else {
				rangeMatch.SetBlackMsg("属性 %v 值 %v 的字符串长度位于禁用的范围 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
			}
		} else if util.IsCheckNumber(field.Type.Kind()) {
			rangeMatch.SetBlackMsg("属性 %v 值 %v 位于禁用的范围 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
		} else if field.Type.Kind() == reflect.Slice {
			if reflect.ValueOf(fieldValue).Len() > 1024 {
				rangeMatch.SetBlackMsg("属性 %v 值的数组长度位于禁用的范围 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
			} else {
				rangeMatch.SetBlackMsg("属性 %v 值 %v 的数组长度位于禁用的范围 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
			}
		} else if field.Type.String() == "time.Time" {
			rangeMatch.SetBlackMsg("属性 %v 时间 %v 位于禁用时间段 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
		} else {
			return true
		}
		return true
	} else {
		if field.Type.Kind() == reflect.String {
			if len(fmt.Sprintf("%v", fieldValue)) > 1024 {
				rangeMatch.SetWhiteMsg("属性 %v 值的长度没有命中只允许的范围 %v", field.Name, fieldValue, rangeMatch.RangeExpress)
			} else {
				rangeMatch.SetWhiteMsg("属性 %v 值 %v 的长度没有命中只允许的范围 %v", field.Name, fieldValue, rangeMatch.RangeExpress)
			}
		} else if util.IsCheckNumber(field.Type.Kind()) {
			rangeMatch.SetWhiteMsg("属性 %v 值 %v 没有命中只允许的范围 %v", field.Name, fieldValue, rangeMatch.RangeExpress)
		} else if field.Type.Kind() == reflect.Slice {
			if reflect.ValueOf(fieldValue).Len() > 1024 {
				rangeMatch.SetWhiteMsg("属性 %v 值的数组长度没有命中只允许的范围 %v", field.Name, fieldValue, rangeMatch.RangeExpress)
			} else {
				rangeMatch.SetWhiteMsg("属性 %v 值 %v 的数组长度没有命中只允许的范围 %v", field.Name, fieldValue, rangeMatch.RangeExpress)
			}
		} else if field.Type.String() == "time.Time" {
			rangeMatch.SetWhiteMsg("属性 %v 时间 %v 没有命中只允许的时间段 %v 中", field.Name, fieldValue, rangeMatch.RangeExpress)
		} else {
			return true
		}
		return false
	}
}

func (rangeMatch *RangeMatch) IsEmpty() bool {
	return rangeMatch.Script == ""
}

func BuildRangeMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constant.MATCH != tagName {
		return
	}

	if !strings.Contains(subCondition, constant.Range) || !strings.Contains(subCondition, constant.EQUAL) {
		return
	}

	index := strings.Index(subCondition, "=")
	value := subCondition[index+1:]

	rangeEntity := parseRange(fieldKind, value)
	if rangeEntity == nil {
		return
	}

	beginAli := rangeEntity.beginAli
	begin := rangeEntity.begin
	end := rangeEntity.end
	endAli := rangeEntity.endAli
	beginNow := rangeEntity.beginNow
	endNow := rangeEntity.endNow

	var script string
	if begin == nil {
		if end == nil {
			if beginNow {
				if constant.LeftEqual == beginAli {
					script = "begin <= value"
				} else if constant.LeftUnEqual == beginAli {
					script = "begin < value"
				}
			} else if endNow {
				if constant.RightEqual == endAli {
					script = "value <= end"
				} else if constant.RightUnEqual == endAli {
					script = "value < end"
				}
			} else {
				return
			}
		} else {
			if beginNow {
				if constant.LeftEqual == beginAli && constant.RightEqual == endAli {
					script = "begin <= value && value <= end"
				} else if constant.LeftEqual == beginAli && constant.RightUnEqual == endAli {
					script = "begin <= value && value < end"
				} else if constant.LeftUnEqual == beginAli && constant.RightEqual == endAli {
					script = "begin < value && value <= end"
				} else if constant.LeftUnEqual == beginAli && constant.RightUnEqual == endAli {
					script = "begin < value && value < end"
				}
			} else {
				if constant.RightEqual == endAli {
					script = "value <= end"
				} else if constant.RightUnEqual == endAli {
					script = "value < end"
				}
			}
		}
	} else {
		if end == nil {
			if endNow {
				if constant.LeftEqual == beginAli && constant.RightEqual == endAli {
					script = "begin <= value && value <= end"
				} else if constant.LeftEqual == beginAli && constant.RightUnEqual == endAli {
					script = "begin <= value && value < end"
				} else if constant.LeftUnEqual == beginAli && constant.RightEqual == endAli {
					script = "begin < value && value <= end"
				} else if constant.LeftUnEqual == beginAli && constant.RightUnEqual == endAli {
					script = "begin < value && value < end"
				}
			} else {
				if constant.LeftEqual == beginAli {
					script = "begin <= value"
				} else if constant.LeftUnEqual == beginAli {
					script = "begin < value"
				}
			}
		} else {
			if constant.LeftEqual == beginAli && constant.RightEqual == endAli {
				script = "begin <= value && value <= end"
			} else if constant.LeftEqual == beginAli && constant.RightUnEqual == endAli {
				script = "begin <= value && value < end"
			} else if constant.LeftUnEqual == beginAli && constant.RightEqual == endAli {
				script = "begin < value && value <= end"
			} else if constant.LeftUnEqual == beginAli && constant.RightUnEqual == endAli {
				script = "begin < value && value < end"
			}
		}
	}

	tree, err := parser.Parse(script)
	if err != nil {
		logger.Error("脚本：%v 解析异常：%v", script, err.Error())
		return
	}

	program, err := compiler.Compile(tree, nil)
	if err != nil {
		logger.Error("脚本: %v 编译异常：%v", script, err.Error())
		return
	}

	addMatcher(objectTypeFullName, objectFieldName, &RangeMatch{Program: program, Begin: begin, End: end, Script: script, RangeExpress: value, BeginNow: beginNow, EndNow: endNow}, errMsg, true)
}

func parseRange(fieldKind reflect.Kind, subCondition string) *RangeEntity {
	subData := rangeRegex.FindAllStringSubmatch(subCondition, -1)
	if len(subData) > 0 {
		beginAli := subData[0][1]
		begin := subData[0][2]
		end := subData[0][5]
		endAli := subData[0][6]

		if (begin == "nil" || begin == "") && (end == "nil" || end == "") {
			logger.Error("range匹配器格式输入错误，start和end不可都为null或者空字符, input=%v", subCondition)
			return nil
		} else if begin == "past" || begin == "future" {
			logger.Error("range匹配器格式输入错误, start不可含有past或者future, input=%v", subCondition)
			return nil
		} else if end == "past" || end == "future" {
			logger.Error("range匹配器格式输入错误, end不可含有past或者future, input=%v", subCondition)
			return nil
		}

		// 如果是数字，则按照数字解析
		if (begin != "" && digitRegex.MatchString(begin)) || (end != "" && digitRegex.MatchString(end)) {
			beginNum := parseNum(fieldKind, begin)
			endNum := parseNum(fieldKind, end)

			return &RangeEntity{beginAli: beginAli, begin: beginNum, end: endNum, endAli: endAli, dateFlag: true}
		} else if (begin != "" && timePlusRegex.MatchString(begin)) || (end != "" && timePlusRegex.MatchString(end)) {
			// 解析动态时间
			dynamicBegin := parseDynamicTime(begin)
			dynamicEnd := parseDynamicTime(end)
			if dynamicBegin == util.EmptyTime && dynamicEnd == util.EmptyTime {
				return nil
			}

			if dynamicBegin == util.EmptyTime {
				return &RangeEntity{beginAli: beginAli, begin: nil, end: dynamicEnd.UnixNano(), endAli: endAli, dateFlag: true}
			} else if dynamicEnd == util.EmptyTime {
				return &RangeEntity{beginAli: beginAli, begin: dynamicBegin.UnixNano(), end: nil, endAli: endAli, dateFlag: true}
			} else {
				return &RangeEntity{beginAli: beginAli, begin: dynamicBegin.UnixNano(), end: dynamicEnd.UnixNano(), endAli: endAli, dateFlag: true}
			}
		} else {
			var beginNow bool
			var endNow bool
			var beginTime time.Time
			var endTime time.Time
			if begin == constant.Now {
				beginNow = true
			} else {
				beginTime = util.ParseTime(begin)
			}

			if end == constant.Now {
				endNow = true
			} else {
				endTime = util.ParseTime(end)
			}

			beginTimeIsEmpty := util.IsTimeEmpty(beginTime)
			endTimeIsEmpty := util.IsTimeEmpty(endTime)

			if !beginTimeIsEmpty && !endTimeIsEmpty {
				if beginTime.After(endTime) {
					logger.Error("时间的范围起始点不正确，起点时间不应该大于终点时间")
					return nil
				}
				return &RangeEntity{beginAli: beginAli, begin: beginTime.UnixNano(), end: endTime.UnixNano(), endAli: endAli, dateFlag: true, beginNow: beginNow, endNow: endNow}
			} else if beginTimeIsEmpty && endTimeIsEmpty {
				logger.Error("range 匹配器格式输入错误，解析数字或者日期失败, time: %v", subData)
			} else {
				if !beginTimeIsEmpty {
					return &RangeEntity{beginAli: beginAli, begin: beginTime.UnixNano(), end: nil, endAli: endAli, dateFlag: true, beginNow: beginNow, endNow: endNow}
				} else if !endTimeIsEmpty {
					return &RangeEntity{beginAli: beginAli, begin: nil, end: endTime.UnixNano(), endAli: endAli, dateFlag: true, beginNow: beginNow, endNow: endNow}
				} else {
					return nil
				}
			}
		}
	} else {
		// 匹配过去和未来的时间
		if subCondition == constant.Past {
			// 过去，则范围为(null, now)
			return &RangeEntity{beginAli: constant.LeftUnEqual, begin: nil, end: nil, endAli: constant.RightUnEqual, dateFlag: true, endNow: true}
		} else if subCondition == constant.Future {
			// 未来，则范围为(now, null)
			return &RangeEntity{beginAli: constant.LeftUnEqual, begin: nil, end: nil, endAli: constant.RightUnEqual, dateFlag: true, beginNow: true}
		}
		return nil
	}
	return nil
}

func parseNum(fieldKind reflect.Kind, valueStr string) interface{} {
	if util.IsCheckNumber(fieldKind) {
		result, err := util.Cast(fieldKind, valueStr)
		if err != nil {
			return nil
		}
		return result
	} else if fieldKind == reflect.String || fieldKind == reflect.Slice {
		result, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil
		}
		return result
	} else {
		return nil
	}
}

func parseDynamicTime(valueStr string) time.Time {
	valueStr = strings.TrimSpace(valueStr)
	if valueStr == "" {
		return util.EmptyTime
	}
	subData := timePlusRegex.FindAllStringSubmatch(valueStr, -1)
	if len(subData) > 0 {
		plusOrMinus := subData[0][1]
		var years, months, days int
		yearStr := subData[0][2]
		monthStr := subData[0][3]
		dayStr := subData[0][4]
		if yearStr != "" {
			yearStr = yearStr[:len(yearStr)-1]
		}
		if monthStr != "" {
			monthStr = monthStr[:len(monthStr)-1]
		}
		if dayStr != "" {
			dayStr = dayStr[:len(dayStr)-1]
		}
		years, _ = strconv.Atoi(fmt.Sprintf("%v%v", plusOrMinus, yearStr))
		months, _ = strconv.Atoi(fmt.Sprintf("%v%v", plusOrMinus, monthStr))
		days, _ = strconv.Atoi(fmt.Sprintf("%v%v", plusOrMinus, dayStr))

		hours := subData[0][5]
		minutes := subData[0][6]
		seconds := subData[0][7]

		resultTime := util.AddYears(time.Now(), years)
		resultTime = util.AddMonths(resultTime, months)
		resultTime = util.AddDays(resultTime, days)
		resultTime = util.AddHour(resultTime, plusOrMinus, hours)
		resultTime = util.AddMinutes(resultTime, plusOrMinus, minutes)
		resultTime = util.AddSeconds(resultTime, plusOrMinus, seconds)

		return resultTime
	}
	return util.EmptyTime
}
