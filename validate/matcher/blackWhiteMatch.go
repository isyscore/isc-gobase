package matcher

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/parser"
	"github.com/isyscore/isc-gobase/logger"
)

type BlackWhiteMatch struct {
	BlackMsg string
	WhiteMsg string
}

func (blackWhiteMatch *BlackWhiteMatch) SetBlackMsg(format string, a ...any) {
	blackWhiteMatch.BlackMsg = fmt.Sprintf(format, a...)
}

func (blackWhiteMatch *BlackWhiteMatch) SetWhiteMsg(format string, a ...any) {
	blackWhiteMatch.WhiteMsg = fmt.Sprintf(format, a...)
}

func (blackWhiteMatch *BlackWhiteMatch) GetWhitMsg() string {
	return blackWhiteMatch.WhiteMsg
}

func (blackWhiteMatch *BlackWhiteMatch) GetBlackMsg() string {
	return blackWhiteMatch.BlackMsg
}

func addMatcher(objectTypeFullName string, objectFieldName string, matcher Matcher, errMsg string, accept bool) {
	fieldMatcherMap, c1 := MatchMap[objectTypeFullName]

	if !c1 {
		fieldMap := make(map[string]*FieldMatcher)
		var matchers []*Matcher
		if matcher != nil {
			matchers = append(matchers, &matcher)
		}

		if errMsg != "" {
			errMsgData := errMsgChange(errMsg)
			tree, err := parser.Parse(errMsgData)
			if err != nil {
				logger.Error("errMsg[%v] parse error: %v", errMsg, err.Error())
				return
			}

			program, err := compiler.Compile(tree, nil)
			if err != nil {
				logger.Error("errMsg[%v] compile error: %v", errMsg, err.Error())
				return
			}

			fieldMap[objectFieldName] = &FieldMatcher{FieldName: objectFieldName, ErrMsgProgram: program, Matchers: matchers, Accept: accept}
		} else {
			fieldMap[objectFieldName] = &FieldMatcher{FieldName: objectFieldName, Matchers: matchers, Accept: accept}
		}

		MatchMap[objectTypeFullName] = fieldMap
	} else {
		fieldMatcher, c2 := fieldMatcherMap[objectFieldName]
		if !c2 {
			var matchers []*Matcher
			if matcher != nil {
				matchers = append(matchers, &matcher)
			}

			if errMsg != "" {
				tree, err := parser.Parse(errMsgChange(errMsg))
				if err != nil {
					logger.Error("errMsg[%v] parse error: %v", errMsg, err.Error())
					return
				}

				program, err := compiler.Compile(tree, nil)
				if err != nil {
					logger.Error("errMsg[%v] compile error: %v", errMsg, err.Error())
					return
				}

				fieldMatcherMap[objectFieldName] = &FieldMatcher{FieldName: objectFieldName, ErrMsgProgram: program, Matchers: matchers, Accept: accept}
			} else {
				fieldMatcherMap[objectFieldName] = &FieldMatcher{FieldName: objectFieldName, Matchers: matchers, Accept: accept}
			}
		} else {
			if matcher != nil {
				fieldMatcher.Matchers = append(fieldMatcher.Matchers, &matcher)
			}
			fieldMatcher.Accept = accept
		}
	}
}

// 将#root和#current转换为root和#current，相当于移除井号
func rmvWell(expression string) string {
	if strings.Contains(expression, "#root.") {
		expression = strings.ReplaceAll(expression, "#root.", "root.")
	}

	if strings.Contains(expression, "#current") {
		expression = strings.ReplaceAll(expression, "#current", "current")
	}
	return expression
}

var currentKey = "#current"
var rootKey = "#root"

func errMsgChange(errMsg string) string {
	var matchKeys []string
	var chgMsg strings.Builder
	chgMsg.WriteString("sprintf(\"")

	var b strings.Builder
	b.Grow(len(errMsg))

	matchIndex := 0
	matchLength := 0
	for infoIndex, data := range errMsg {
		c := string(data)
		if c == "#" {
			if findCurrentKey(infoIndex, 0, errMsg) {
				matchIndex = 0
				matchLength = len(currentKey)
				b.WriteString("%v")
				matchKeys = append(matchKeys, "current")
				continue
			} else if find, size, wordKey := findRootKey(infoIndex, 0, errMsg); find {
				matchIndex = 0
				matchLength = size
				b.WriteString("%v")
				matchKeys = append(matchKeys, "root"+wordKey)
				continue
			}
		} else if matchIndex+1 < matchLength {
			matchIndex++
			continue
		} else {
			b.WriteString(c)
		}
	}

	chgMsg.WriteString(b.String())
	chgMsg.WriteString("\"")

	matchKeysSize := len(matchKeys)
	if matchKeysSize > 0 {
		chgMsg.WriteString(", ")
	}

	for i, data := range matchKeys {
		if i+1 < matchKeysSize {
			chgMsg.WriteString(data)
			chgMsg.WriteString(", ")
		} else {
			chgMsg.WriteString(data)
		}
	}
	chgMsg.WriteString(")")

	return chgMsg.String()
}

func findCurrentKey(infoIndex, matchIndex int, info string) bool {
	if matchIndex >= len(currentKey) {
		return true
	}
	if info[infoIndex:infoIndex+1] == currentKey[matchIndex:matchIndex+1] {
		return findCurrentKey(infoIndex+1, matchIndex+1, info)
	}
	return false
}

func findRootKey(infoIndex, matchIndex int, info string) (bool, int, string) {
	if matchIndex >= len(rootKey) {
		nextKeyLength := nextMatchKeyLength(info[infoIndex:])
		if nextKeyLength > 0 {
			return true, len(rootKey) + nextKeyLength, info[infoIndex : infoIndex+nextKeyLength]
		}
		return false, 0, ""
	}
	if info[infoIndex:infoIndex+1] == rootKey[matchIndex:matchIndex+1] {
		return findRootKey(infoIndex+1, matchIndex+1, info)
	}
	return false, 0, ""
}

// 下一个英文的单词长度
// 97 ~ 122
// 65 ~ 90
func nextMatchKeyLength(errMsg string) int {
	spaceIndex := strings.Index(strings.TrimSpace(errMsg), " ")
	toMatchMsg := errMsg
	if spaceIndex > 0 {
		toMatchMsg = errMsg[:spaceIndex]
	}
	var index = 0
	for _, c := range toMatchMsg {
		// 判断是否是英文字符：a~z、A~Z和点号"."
		if (c >= 97 && c <= 122) || (c >= 65 && c <= 90) || c == 46 {
			index++
			continue
		} else {
			return index
		}
	}
	return index
}
