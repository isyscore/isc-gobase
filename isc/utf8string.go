package isc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ISCUTF8String []rune

func NewUTF8String(str string) ISCUTF8String {
	return ISCUTF8String(str)
}

func (s ISCUTF8String) String() string {
	return string(s)
}

func (s ISCUTF8String) At(index int) rune {
	return s[index]
}

func (s ISCUTF8String) Length() int {
	return len(s)
}

func (s ISCUTF8String) Chars() ISCList[rune] {
	return ISCList[rune](s)
}

func (s ISCUTF8String) Count(substr ISCUTF8String) int {
	return strings.Count(string(s), string(substr))
}

func (s ISCUTF8String) Contains(substr ISCUTF8String) bool {
	return strings.Contains(string(s), string(substr))
}

func (s ISCUTF8String) ContainsRune(r rune) bool {
	return strings.ContainsRune(string(s), r)
}

func (s ISCUTF8String) toIndex(idx int) int {
	if idx < 0 {
		return idx
	}
	buf := string(s)[:idx]
	runeBuf := []rune(buf)
	return len(runeBuf)
}

func (s ISCUTF8String) IndexOf(substr ISCUTF8String) int {
	i := strings.Index(string(s), string(substr))
	return s.toIndex(i)
}

func (s ISCUTF8String) LastIndexOf(substr ISCUTF8String) int {
	i := strings.LastIndex(string(s), string(substr))
	return s.toIndex(i)
}

func (s ISCUTF8String) IndexOfFunc(f func(rune) bool) int {
	i := strings.IndexFunc(string(s), f)
	return s.toIndex(i)
}

func (s ISCUTF8String) LastIndexOfFunc(f func(rune) bool) int {
	i := strings.LastIndexFunc(string(s), f)
	return s.toIndex(i)
}

func (s ISCUTF8String) SplitN(sep ISCUTF8String, n int) []ISCUTF8String {
	ss := strings.SplitN(string(s), string(sep), n)
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) SplitAfterN(sep ISCUTF8String, n int) []ISCUTF8String {
	ss := strings.SplitAfterN(string(s), string(sep), n)
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) Split(sep ISCUTF8String) []ISCUTF8String {
	ss := strings.Split(string(s), string(sep))
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) SplitAfter(sep ISCUTF8String) []ISCUTF8String {
	ss := strings.SplitAfter(string(s), string(sep))
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) Fields() []ISCUTF8String {
	ss := strings.Fields(string(s))
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) FieldsFunc(f func(rune) bool) []ISCUTF8String {
	ss := strings.FieldsFunc(string(s), f)
	return ListToMapFrom[string, ISCUTF8String](ss).Map(func(item string) ISCUTF8String {
		return ISCUTF8String(item)
	})
}

func (s ISCUTF8String) StartsWith(prefix ISCUTF8String) bool {
	return strings.HasPrefix(string(s), string(prefix))
}

func (s ISCUTF8String) EndsWith(suffix ISCUTF8String) bool {
	return strings.HasSuffix(string(s), string(suffix))
}

func (s ISCUTF8String) TrimLeftFunc(f func(rune) bool) ISCUTF8String {
	return ISCUTF8String(strings.TrimLeftFunc(string(s), f))
}

func (s ISCUTF8String) TrimRightFunc(f func(rune) bool) ISCUTF8String {
	return ISCUTF8String(strings.TrimRightFunc(string(s), f))
}

func (s ISCUTF8String) TrimFunc(f func(rune) bool) ISCUTF8String {
	return ISCUTF8String(strings.TrimFunc(string(s), f))
}

func (s ISCUTF8String) Trim(cutset ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.Trim(string(s), string(cutset)))
}

func (s ISCUTF8String) TrimLeft(cutset ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.TrimLeft(string(s), string(cutset)))
}

func (s ISCUTF8String) TrimRight(cutset ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.TrimRight(string(s), string(cutset)))
}

func (s ISCUTF8String) TrimSpace() ISCUTF8String {
	return ISCUTF8String(strings.TrimSpace(string(s)))
}

func (s ISCUTF8String) TrimPrefix(prefix ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.TrimPrefix(string(s), string(prefix)))
}

func (s ISCUTF8String) TrimSuffix(suffix ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.TrimSuffix(string(s), string(suffix)))
}

func (s ISCUTF8String) Replace(old, new ISCUTF8String, n int) ISCUTF8String {
	return ISCUTF8String(strings.Replace(string(s), string(old), string(new), n))
}

func (s ISCUTF8String) ReplaceAll(old, new ISCUTF8String) ISCUTF8String {
	return ISCUTF8String(strings.ReplaceAll(string(s), string(old), string(new)))
}

func (s ISCUTF8String) EqualFold(t ISCUTF8String) bool {
	return strings.EqualFold(string(s), string(t))
}

func (s ISCUTF8String) ToUpper() ISCUTF8String {
	return ISCUTF8String(strings.ToUpper(string(s)))
}

func (s ISCUTF8String) ToLower() ISCUTF8String {
	return ISCUTF8String(strings.ToLower(string(s)))
}

func (s ISCUTF8String) ToTitle() ISCUTF8String {
	return ISCUTF8String(strings.ToTitle(string(s)))
}

func (s ISCUTF8String) IsEmpty() bool {
	return len(s) == 0
}

func (s ISCUTF8String) SubStringStart(AStartIndex int) ISCUTF8String {
	return s[AStartIndex:]
}

func (s ISCUTF8String) SubStringStartEnd(AStartIndex, AEndIndex int) ISCUTF8String {
	return s[AStartIndex:AEndIndex]
}

func (s ISCUTF8String) SubStringBefore(delimiter ISCUTF8String) ISCUTF8String {
	if i := s.IndexOf(delimiter); i != -1 {
		return s[:i]
	} else {
		return s
	}
}

func (s ISCUTF8String) SubStringAfter(delimiter ISCUTF8String) ISCUTF8String {
	if i := s.IndexOf(delimiter); i != -1 {
		return s[i+len(delimiter):]
	} else {
		return s
	}
}

func (s ISCUTF8String) SubStringBeforeLast(delimiter ISCUTF8String) ISCUTF8String {
	if i := s.LastIndexOf(delimiter); i != -1 {
		return s[:i]
	} else {
		return s
	}
}

func (s ISCUTF8String) SubStringAfterLast(delimiter ISCUTF8String) ISCUTF8String {
	if i := s.LastIndexOf(delimiter); i != -1 {
		return s[i+len(delimiter):]
	} else {
		return s
	}
}

func (s ISCUTF8String) Insert(index int, substr ISCUTF8String) ISCUTF8String {
	ss := string(s[:index]) + string(substr) + string(s[index:])
	return ISCUTF8String(ss)
}

func (s ISCUTF8String) Delete(index int, count int) ISCUTF8String {
	ss := string(s[:index]) + string(s[index+count:])
	return ISCUTF8String(ss)
}

func (s ISCUTF8String) Matches(pattern ISCUTF8String) bool {
	reg := regexp.MustCompile(string(pattern))
	return reg.MatchString(string(s))
}

func (s ISCUTF8String) Lines() []ISCUTF8String {
	return s.Split(ISCUTF8String("\n"))
}

func (s ISCUTF8String) LinesNoEmpty() []ISCUTF8String {
	return NewListWithList(s.Split(ISCUTF8String("\n"))).Filter(func(item ISCUTF8String) bool {
		return !item.IsEmpty()
	})
}

func (s ISCUTF8String) ToBoolean() bool {
	return string(s.ToLower()) == "true"
}

func (s ISCUTF8String) ToInt() int {
	return ToInt(string(s))
}

func (s ISCUTF8String) ToInt8() int8 {
	return ToInt8(string(s))
}

func (s ISCUTF8String) ToInt16() int16 {
	return ToInt16(string(s))
}

func (s ISCUTF8String) ToInt32() int32 {
	return ToInt32(string(s))
}

func (s ISCUTF8String) ToInt64() int64 {
	return ToInt64(string(s))
}

func (s ISCUTF8String) ToFloat() float32 {
	return ToFloat32(string(s))
}

func (s ISCUTF8String) ToFloat64() float64 {
	return ToFloat64(string(s))
}

func (s ISCUTF8String) ToIntRadix(radix int) (int64, error) {
	if radix != 2 && radix != 8 && radix != 10 && radix != 16 {
		return 0, fmt.Errorf("radix %d is not supported", radix)
	}
	size := strconv.IntSize
	return strconv.ParseInt(string(s), radix, size)
}

func (s ISCUTF8String) ToJSONEncoded() ISCUTF8String {
	return s.ReplaceAll(ISCUTF8String("\\"), ISCUTF8String("\\\\")).ReplaceAll(ISCUTF8String("\n"), ISCUTF8String("\\n")).ReplaceAll(ISCUTF8String("\""), ISCUTF8String("\\\""))
}

func (s ISCUTF8String) ToPair() Pair[ISCUTF8String, ISCUTF8String] {
	sa := s.TrimSpace().Split(ISCUTF8String("="))
	return NewPair(sa[0].TrimSpace(), sa[1].TrimSpace())
}
