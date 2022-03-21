package matcher

import (
	"reflect"

	"github.com/antonmedv/expr/vm"
)

type Matcher interface {
	Match(object any, field reflect.StructField, fieldValue any) bool
	IsEmpty() bool
	GetWhitMsg() string
	GetBlackMsg() string
}

type FieldMatcher struct {

	// 属性名
	FieldName string
	// 异常信息编译后的处理
	ErrMsgProgram *vm.Program
	// 是否接受：true，则表示白名单，false，则表示黑名单
	Accept bool
	// 是否禁用
	Disable bool
	// 匹配器列表
	Matchers []*Matcher
}

// MatchMap key：类全名，value：key：属性名
var MatchMap = make(map[string]map[string]*FieldMatcher)
