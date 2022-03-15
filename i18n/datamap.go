package i18n

import (
	f0 "github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/isc"
)

type I18NMap struct {
	Language        string            // 当前的语言
	DefaultLanguage string            // 默认语言
	Data            map[string]string // 当前的字符串映射
	DefaultData     map[string]string // 默认的字符串映射，当从Data内找不到key时，从此处找，再找不到就报错
}

var innerMap *I18NMap

func NewI18NMap(language string, filePath string) *I18NMap {
	m := loadPo(filePath)
	return &I18NMap{
		Language:        language,
		DefaultLanguage: language,
		Data:            m,
		DefaultData:     m,
	}
}

// po 文件格式
// msgid "hello"
// msgstr "你好"

func loadPo(filePath string) map[string]string {
	m := make(map[string]string)
	lines := f0.ReadFileLines(filePath)
	for _, s := range lines {
		ss := isc.ISCString(s)
		key := ss.SubStringBefore(" ")
		val := ss.SubStringAfter(" ").TrimSpace().Trim("\"").ReplaceAll("\\n", "\n").ReplaceAll("\\r", "\r").ReplaceAll("\\t", "\t").ReplaceAll("\\\"", "\"").ReplaceAll("\\\\", "\\")
		m[string(key)] = string(val)
	}
	return m
}
