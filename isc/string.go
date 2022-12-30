package isc

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"regexp"
	"strconv"
	"strings"
)

type ISCString string

func (s ISCString) At(index int) uint8 {
	return s[index]
}

func (s ISCString) Length() int {
	return len(s)
}

func (s ISCString) Chars() ISCList[uint8] {
	var list []uint8
	for i := 0; i < len(s); i++ {
		list = append(list, s[i])
	}
	return list
}

func (s ISCString) Count(substr string) int {
	return strings.Count(string(s), substr)
}

func (s ISCString) Contains(substr string) bool {
	return strings.Contains(string(s), substr)
}

func (s ISCString) ContainsAny(chars string) bool {
	return strings.ContainsAny(string(s), chars)
}

func (s ISCString) ContainsRune(r rune) bool {
	return strings.ContainsRune(string(s), r)
}

func (s ISCString) IndexOf(substr string) int {
	return strings.Index(string(s), substr)
}

func (s ISCString) LastIndexOf(substr string) int {
	return strings.LastIndex(string(s), substr)
}

func (s ISCString) IndexByteOf(c byte) int {
	return strings.IndexByte(string(s), c)
}

func (s ISCString) IndexOfAny(chars string) int {
	return strings.IndexAny(string(s), chars)
}

func (s ISCString) LastIndexOfAny(chars string) int {
	return strings.LastIndexAny(string(s), chars)
}

func (s ISCString) LastIndexOfByte(c byte) int {
	return strings.LastIndexByte(string(s), c)
}

func (s ISCString) SplitN(sep string, n int) []ISCString {
	ss := strings.SplitN(string(s), sep, n)
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) SplitAfterN(sep string, n int) []ISCString {
	ss := strings.SplitAfterN(string(s), sep, n)
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) Split(sep string) []ISCString {
	ss := strings.Split(string(s), sep)
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) SplitAfter(sep string) []ISCString {
	ss := strings.SplitAfter(string(s), sep)
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) Fields() []ISCString {
	ss := strings.Fields(string(s))
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) FieldsFunc(f func(rune) bool) []ISCString {
	ss := strings.FieldsFunc(string(s), f)
	return ListToMapFrom[string, ISCString](ss).Map(func(item string) ISCString {
		return ISCString(item)
	})
}

func (s ISCString) StartsWith(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

func (s ISCString) EndsWith(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

func (s ISCString) Repeat(count int) ISCString {
	return ISCString(strings.Repeat(string(s), count))
}

func (s ISCString) TrimLeftFunc(f func(rune) bool) ISCString {
	return ISCString(strings.TrimLeftFunc(string(s), f))
}

func (s ISCString) TrimRightFunc(f func(rune) bool) ISCString {
	return ISCString(strings.TrimRightFunc(string(s), f))
}

func (s ISCString) TrimFunc(f func(rune) bool) ISCString {
	return ISCString(strings.TrimFunc(string(s), f))
}

func (s ISCString) IndexOfFunc(f func(rune) bool) int {
	return strings.IndexFunc(string(s), f)
}

func (s ISCString) LastIndexOfFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(string(s), f)
}

func (s ISCString) Trim(cutset string) ISCString {
	return ISCString(strings.Trim(string(s), cutset))
}

func (s ISCString) TrimLeft(cutset string) ISCString {
	return ISCString(strings.TrimLeft(string(s), cutset))
}

func (s ISCString) TrimRight(cutset string) ISCString {
	return ISCString(strings.TrimRight(string(s), cutset))
}

func (s ISCString) TrimSpace() ISCString {
	return ISCString(strings.TrimSpace(string(s)))
}

func (s ISCString) TrimPrefix(prefix string) ISCString {
	return ISCString(strings.TrimPrefix(string(s), prefix))
}

func (s ISCString) TrimSuffix(suffix string) ISCString {
	return ISCString(strings.TrimSuffix(string(s), suffix))
}

func (s ISCString) Replace(old, new string, n int) ISCString {
	return ISCString(strings.Replace(string(s), old, new, n))
}

func (s ISCString) ReplaceAll(old, new string) ISCString {
	return ISCString(strings.ReplaceAll(string(s), old, new))
}

func (s ISCString) EqualFold(t string) bool {
	return strings.EqualFold(string(s), t)
}

func (s ISCString) ToUpper() ISCString {
	return ISCString(strings.ToUpper(string(s)))
}

func (s ISCString) ToLower() ISCString {
	return ISCString(strings.ToLower(string(s)))
}

func (s ISCString) ToTitle() ISCString {
	return ISCString(strings.ToTitle(string(s)))
}

func (s ISCString) IsEmpty() bool {
	return len(s) == 0
}

func (s ISCString) SubStringStart(AStartIndex int) ISCString {
	return s[AStartIndex:]
}

func (s ISCString) SubStringStartEnd(AStartIndex, AEndIndex int) ISCString {
	return s[AStartIndex:AEndIndex]
}

func (s ISCString) SubStringBefore(delimiter string) ISCString {
	if i := s.IndexOf(delimiter); i != -1 {
		return s[:i]
	} else {
		return s
	}
}

func (s ISCString) SubStringAfter(delimiter string) ISCString {
	if i := s.IndexOf(delimiter); i != -1 {
		return s[i+len(delimiter):]
	} else {
		return s
	}
}

func (s ISCString) SubStringBeforeLast(delimiter string) ISCString {
	if i := s.LastIndexOf(delimiter); i != -1 {
		return s[:i]
	} else {
		return s
	}
}

func (s ISCString) SubStringAfterLast(delimiter string) ISCString {
	if i := s.LastIndexOf(delimiter); i != -1 {
		return s[i+len(delimiter):]
	} else {
		return s
	}
}

func (s ISCString) Insert(index int, substr string) ISCString {
	ss := string(s[:index]) + substr + string(s[index:])
	return ISCString(ss)
}

func (s ISCString) Delete(index int, count int) ISCString {
	ss := string(s[:index]) + string(s[index+count:])
	return ISCString(ss)
}

func (s ISCString) Matches(pattern string) bool {
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(string(s))
}

func (s ISCString) Lines() []ISCString {
	return s.Split("\n")
}

func (s ISCString) LinesNoEmpty() []ISCString {
	return NewListWithList(s.Split("\n")).Filter(func(item ISCString) bool {
		return !item.IsEmpty()
	})
}

func (s ISCString) ToBoolean() bool {
	return s.ToLower() == "true"
}

func (s ISCString) ToInt() int {
	return ToInt(s)
}

func (s ISCString) ToInt8() int8 {
	return ToInt8(s)
}

func (s ISCString) ToInt16() int16 {
	return ToInt16(s)
}

func (s ISCString) ToInt32() int32 {
	return ToInt32(s)
}

func (s ISCString) ToInt64() int64 {
	return ToInt64(s)
}

func (s ISCString) ToFloat() float32 {
	return ToFloat32(s)
}

func (s ISCString) ToFloat64() float64 {
	return ToFloat64(s)
}

func (s ISCString) ToIntRadix(radix int) (int64, error) {
	if radix != 2 && radix != 8 && radix != 10 && radix != 16 {
		return 0, fmt.Errorf("radix %d is not supported", radix)
	}
	size := strconv.IntSize
	return strconv.ParseInt(string(s), radix, size)
}

func (s ISCString) ToJSONEncoded() ISCString {
	return s.ReplaceAll("\\", "\\\\").ReplaceAll("\n", "\\n").ReplaceAll("\"", "\\\"")
}

func (s ISCString) ToMap() ISCMap[ISCString, ISCString] {
	return ListToTripleFrom[ISCString, ISCString, ISCString](s.Split("&")).Associate(func(item ISCString) Pair[ISCString, ISCString] {
		sa := item.Split("=")
		return NewPair(sa[0], sa[1])
	})
}

func (s ISCString) ToCookieMap() ISCMap[ISCString, ISCString] {
	return ListToTripleFrom[ISCString, ISCString, ISCString](s.Split(";")).Associate(func(item ISCString) Pair[ISCString, ISCString] {
		sa := item.TrimSpace().Split("=")
		return NewPair(sa[0].TrimSpace(), sa[1].TrimSpace())
	})
}

func (s ISCString) ToPair() Pair[ISCString, ISCString] {
	sa := s.TrimSpace().Split("=")
	return NewPair(sa[0].TrimSpace(), sa[1].TrimSpace())
}

func (s ISCString) Drop(n int) ISCString {
	return s[n:]
}

func (s ISCString) DropLast(n int) ISCString {
	return s[:len(s)-n]
}

func (s ISCString) Take(n int) ISCString {
	return s[:n]
}

func (s ISCString) TakeLast(n int) ISCString {
	return s[len(s)-n:]
}

// BigCamel 小驼峰到大驼峰：首字母变成大写: dataBaseUser -> DateBaseUser
func BigCamel(word string) string {
	if word == "" {
		return ""
	}
	return strings.ToUpper(word[:1]) + word[1:]
}

// BigCamelToMiddleLine 大驼峰到中划线: DataBaseUser -> data-db-user
func BigCamelToMiddleLine(word string) string {
	if word == "" {
		return ""
	}
	return MiddleLine(BigCamelToSmallCamel(word))
}

// BigCamelToPostUnder 大驼峰到后缀下划线: DataBaseUser -> data_base_user_
func BigCamelToPostUnder(word string) string {
	if word == "" {
		return ""
	}
	return PostUnder(BigCamelToSmallCamel(word))
}

// BigCamelToPrePostUnder 大驼峰到前后缀下划线: DataBaseUser -> _data_base_user_
func BigCamelToPrePostUnder(word string) string {
	if word == "" {
		return ""
	}
	return PrePostUnder(BigCamelToSmallCamel(word))
}

// BigCamelToPreUnder 大驼峰到前后缀下划线: DataBaseUser -> _data_base_user
func BigCamelToPreUnder(word string) string {
	if word == "" {
		return ""
	}
	return PreUnder(BigCamelToSmallCamel(word))
}

// BigCamelToSmallCamel 大驼峰到小驼峰：首字母变成小写：DataBaseUser -> dataBaseUser
func BigCamelToSmallCamel(word string) string {
	if word == "" {
		return ""
	}
	return strings.ToLower(word[:1]) + word[1:]
}

// BigCamelToUnderLine 大驼峰到下划线：DataBaseUser -> data_base_user
func BigCamelToUnderLine(word string) string {
	if word == "" {
		return ""
	}
	return UnderLine(BigCamelToSmallCamel(word))
}

// BigCamelToUpperMiddle 大驼峰到小写中划线：DataBaseUser -> DATA-BASE-USER
func BigCamelToUpperMiddle(word string) string {
	if word == "" {
		return ""
	}
	return UpperUnderMiddle(BigCamelToSmallCamel(word))
}

// BigCamelToUpperUnder 大驼峰到大写下划线: DataBaseUser -> DATA_BASE_USER
func BigCamelToUpperUnder(word string) string {
	if word == "" {
		return ""
	}
	return UpperUnder(BigCamelToSmallCamel(word))
}

// MiddleLine 小驼峰到中划线：dataBaseUser -> data-db-user
func MiddleLine(word string) string {
	if word == "" {
		return ""
	}
	reg, err := regexp.Compile("\\B[A-Z]")
	if err != nil {
		return word
	}

	subIndex := reg.FindAllStringSubmatchIndex(word, -1)
	var lastIndex = 0
	var result = ""
	for i := 0; i < len(subIndex); i++ {
		result += word[lastIndex:subIndex[i][0]]
		result += "-" + strings.ToLower(word[subIndex[i][0]:subIndex[i][1]])
		lastIndex = subIndex[i][1]
	}
	result += word[lastIndex:]
	return result
}

// MiddleLineToBigCamel 中划线到大驼峰：data-db-user -> DataBaseUser
func MiddleLineToBigCamel(word string) string {
	if word == "" {
		return ""
	}
	return BigCamel(MiddleLineToSmallCamel(word))
}

// MiddleLineToSmallCamel 中划线到小驼峰：data-base-user -> dataBaseUser
func MiddleLineToSmallCamel(word string) string {
	if word == "" {
		return ""
	}
	return strings.ReplaceAll(ToUpperWord("(?<=-)[a-z]", word), "-", "")
}

// PostUnder 小驼峰到后下划线：dataBaseUser -> data_base_user_
func PostUnder(word string) string {
	if word == "" {
		return ""
	}
	return UnderLine(word) + "_"
}

// PreFixUnderLine 小驼峰到添加前缀字符下划线：dataBaseUser -> pre_data_base_user
func PreFixUnderLine(word, preFix string) string {
	return preFix + UnderLine(word)
}

// PreFixUnderToSmallCamel 前缀字符下划线去掉到小驼峰：pre_data_base_user -> dataBaseUser
func PreFixUnderToSmallCamel(word, preFix string) string {
	if strings.HasPrefix(word, preFix) {
		return UnderLineToSmallCamel(word[len(preFix):])
	}
	return UnderLineToSmallCamel(word)
}

// PrePostUnder 小驼峰到前后缀下划线：dataBaseUser -> _data_base_user_
func PrePostUnder(word string) string {
	if word == "" {
		return ""
	}
	return "_" + UnderLine(word) + "_"
}

// PreUnder 小驼峰到前下划线：dataBaseUser -> _data_base_user
func PreUnder(word string) string {
	if word == "" {
		return ""
	}
	return "_" + UnderLine(word)
}

// UnderLine 小驼峰到下划线：非边缘单词开头大写变前下划线和后面大写：dataBaseUser -> data_base_user
func UnderLine(word string) string {
	if word == "" {
		return ""
	}
	reg, err := regexp.Compile("\\B[A-Z]")
	if err != nil {
		return word
	}

	subIndex := reg.FindAllStringSubmatchIndex(word, -1)
	var lastIndex = 0
	var result = ""
	for i := 0; i < len(subIndex); i++ {
		result += word[lastIndex:subIndex[i][0]]
		result += "_" + strings.ToLower(word[subIndex[i][0]:subIndex[i][1]])
		lastIndex = subIndex[i][1]
	}
	result += word[lastIndex:]
	return result
}

// UnderLineToBigCamel 下划线到大驼峰：下划线后面小写变大写，下划线去掉
// data_base_user   -> DataBaseUser
// _data_base_user  -> DataBaseUser
// _data_base_user_ -> DataBaseUser
// data_base_user_  -> DataBaseUser
func UnderLineToBigCamel(word string) string {
	if word == "" {
		return ""
	}
	return BigCamel(strings.ReplaceAll(ToUpperWord("(?<=_)[a-z]", word), "_", ""))
}

// UnderLineToSmallCamel 下划线到小驼峰：下划线后面小写变大写，下划线去掉
// data_base_user   -> dataBaseUser
// _data_base_user  -> dataBaseUser
// _data_base_user_ -> dataBaseUser
// data_base_user_  -> dataBaseUser
func UnderLineToSmallCamel(word string) string {
	if word == "" {
		return ""
	}
	return BigCamelToSmallCamel(strings.ReplaceAll(ToUpperWord("(?<=_)[a-z]", word), "_", ""))
}

// 大写中划线到大驼峰：DATA-BASE-USER -> DataBaseUser
func UpperMiddleToBigCamel(word string) string {
	return BigCamel(UpperUnderMiddleToSmallCamel(word))
}

// UpperUnder 小驼峰到大写下划线：dataBaseUser -> DATA_BASE_USER
func UpperUnder(word string) string {
	if word == "" {
		return ""
	}
	reg, err := regexp.Compile("\\B[A-Z]")
	if err != nil {
		return word
	}

	subIndex := reg.FindAllStringSubmatchIndex(word, -1)
	var lastIndex = 0
	var result = ""
	for i := 0; i < len(subIndex); i++ {
		result += strings.ToUpper(word[lastIndex:subIndex[i][0]])
		result += "_" + strings.ToUpper(word[subIndex[i][0]:subIndex[i][1]])
		lastIndex = subIndex[i][1]
	}
	result += strings.ToUpper(word[lastIndex:])
	return result
}

// UpperUnderMiddle 小驼峰到大写中划线：dataBaseUser -> DATA-BASE-USER
func UpperUnderMiddle(word string) string {
	if word == "" {
		return ""
	}
	reg, err := regexp.Compile("\\B[A-Z]")
	if err != nil {
		return word
	}

	subIndex := reg.FindAllStringSubmatchIndex(word, -1)
	var lastIndex = 0
	var result = ""
	for i := 0; i < len(subIndex); i++ {
		result += strings.ToUpper(word[lastIndex:subIndex[i][0]])
		result += "-" + strings.ToUpper(word[subIndex[i][0]:subIndex[i][1]])
		lastIndex = subIndex[i][1]
	}
	result += strings.ToUpper(word[lastIndex:])
	return result
}

// UpperUnderMiddleToSmallCamel 大写中划线到大驼峰：DATA-BASE-USER -> dataBaseUser
func UpperUnderMiddleToSmallCamel(word string) string {
	if word == "" {
		return ""
	}
	return MiddleLineToSmallCamel(strings.ToLower(word))
}

// UpperUnderToBigCamel 大写下划线到大驼峰：DATA_BASE_USER -> DataBaseUser
func UpperUnderToBigCamel(word string) string {
	if word == "" {
		return ""
	}
	return BigCamel(UpperUnderToSmallCamel(word))
}

// UpperUnderToSmallCamel 大写下划线到小驼峰：DATA_BASE_USER -> dataBaseUser
func UpperUnderToSmallCamel(word string) string {
	if word == "" {
		return ""
	}
	return UnderLineToSmallCamel(strings.ToLower(word))
}

// ToUpperWord 匹配的单词变为大写
// regex: 正则表达式，主要用于匹配某些字符变为大写
// word: 待匹配字段
func ToUpperWord(regex, word string) string {
	if word == "" {
		return ""
	}
	regexResult, err := regexp2.Compile(regex, 0)
	if err != nil {
		return word
	}

	matcherResult, err := regexResult.FindStringMatch(word)
	if err != nil {
		return word
	}

	var result = ""
	var lastIndex = 0
	for matcherResult != nil {
		result += word[lastIndex:matcherResult.Index]
		result += strings.ToUpper(word[matcherResult.Index : matcherResult.Index+1])
		lastIndex = matcherResult.Index + 1
		matcherResult, err = regexResult.FindNextMatch(matcherResult)
		if err != nil {
			continue
		}
	}
	result += word[lastIndex:]
	return result
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
	PB = 1024 * TB
	EB = 1024 * PB
	//ZB = 1024*EB
	//YB = 1024*ZB
	//BB = 1024*YB
)

func FormatSize(fileSize int64) (size string) {
	if fileSize < KB {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(B))
	} else if fileSize < MB {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(KB))
	} else if fileSize < GB {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(MB))
	} else if fileSize < TB {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(GB))
	} else if fileSize < PB {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(TB))
	} else if fileSize < EB {
		return fmt.Sprintf("%.2fPB", float64(fileSize)/float64(PB))
	} else {
		// 不要加更多判断了，编译器报错
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(EB))
	}
}

func ParseByteSize(byteStr string) (size int64) {
	if byteStr == "" {
		return 0
	}

	if strings.HasSuffix(byteStr, "GB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * GB
	}

	if strings.HasSuffix(byteStr, "MB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * MB
	}

	if strings.HasSuffix(byteStr, "KB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * KB
	}

	byteStr = strings.ToUpper(byteStr)
	if strings.HasSuffix(byteStr, "EB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * EB
	}

	if strings.HasSuffix(byteStr, "PB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * PB
	}

	if strings.HasSuffix(byteStr, "TB") {
		byteNum := ToInt64(byteStr[:len(byteStr)-2])
		return byteNum * TB
	}

	if strings.HasSuffix(byteStr, "B") {
		byteNum := ToInt64(byteStr[:len(byteStr)-1])
		return byteNum * B
	}

	// 其他的暂时不支持
	return 0
}
