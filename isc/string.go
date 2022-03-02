package isc

import (
	"regexp"
	"strings"
)

type ISCString string

func (s ISCString) At(index int) uint8 {
	return s[index]
}

func (s ISCString) Length() int {
	return len(s)
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

func (s ISCString) IndexOf(substr string) int {
	return strings.Index(string(s), substr)
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

func (s ISCString) ToInt64() int64 {
	return ToInt64(s)
}

func (s ISCString) ToFloat() float32 {
	return ToFloat32(s)
}

func (s ISCString) ToFloat64() float64 {
	return ToFloat64(s)
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
