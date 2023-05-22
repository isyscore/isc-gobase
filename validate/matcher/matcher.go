package matcher

import (
	"github.com/antonmedv/expr/vm"
	"reflect"
)

type Matcher interface {
	//Match(object any, field reflect.StructField, fieldValue any) bool
	Match(parameterMap map[string]interface{}, object interface{}, field reflect.StructField, fieldValue interface{}) bool
	IsEmpty() bool
	GetErrCode() string
	GetWhitMsg() string
	GetBlackMsg() string
}

type FieldMatcher struct {

	// 属性名
	FieldName string
	// 匹配异常后的code
	ErrCode string
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
