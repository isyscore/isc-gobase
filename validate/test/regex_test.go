package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type ValueRegexEntity struct {
	Name string `match:"regex=^zhou.*zhen$"`
	Age  int    `match:"regex=^\\d+$"`
}

func TestRegex(t *testing.T) {
	var value ValueRegexEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueRegexEntity{Name: "zhouOKzhen"}
	result, err = validate.Check(value, "name")
	TrueErr(t, result, err)

	// 测试 正常情况
	value = ValueRegexEntity{Age: 13}
	result, err = validate.Check(value, "age")
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueRegexEntity{Name: "chenzhen"}
	result, err = validate.Check(value, "name")
	Equal(t, "属性 Name 的值 chenzhen 没命中只允许的正则表达式 ^zhou.*zhen$ ", err, result, false)
}

// Regex的基准测试
func Benchmark_Regex(b *testing.B) {
	var value ValueRegexEntity
	for i := 0; i < b.N; i++ {
		value = ValueRegexEntity{Name: "chenzhen"}
		validate.Check(value, "name")
	}
}
