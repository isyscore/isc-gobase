package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type IsBlankEntity1 struct {
	Name string `match:"isBlank=false"`
	Age  int
}

// 默认为true
type IsBlankEntity2 struct {
	Name string `match:"isBlank"`
	Age  int
}

type IsBlankEntity3 struct {
	Name string `match:"isBlank=true"`
	Age  int
}

// 测试基本类型
func TestIsBlank1(t *testing.T) {
	var value IsBlankEntity1
	var result bool
	var err string

	//测试 正常情况
	value = IsBlankEntity1{Name: "zhou"}
	result, err = validate.Check(value, "name")
	TrueErr(t, result, err)

	// 测试 正常情况
	value = IsBlankEntity1{Age: 13}
	result, err = validate.Check(value, "name")
	Equal(t, err, "属性 Name 的值为空字符", result, false)
}

// 测试基本类型：简化版
func TestIsBlank2(t *testing.T) {
	var value IsBlankEntity2
	var result bool
	var err string

	//测试 正常情况
	value = IsBlankEntity2{Name: ""}
	result, err = validate.Check(value, "name")
	TrueErr(t, result, err)

	// 测试 正常情况
	value = IsBlankEntity2{Name: "zhou"}
	result, err = validate.Check(value, "name")
	Equal(t, err, "属性 Name 的值为非空字符", result, false)
}

func TestIsBlank3(t *testing.T) {
	var value IsBlankEntity3
	var result bool
	var err string

	//测试 正常情况
	value = IsBlankEntity3{Name: ""}
	result, err = validate.Check(value, "name")
	TrueErr(t, result, err)

	// 测试 正常情况
	value = IsBlankEntity3{Name: "zhou"}
	result, err = validate.Check(value, "name")
	Equal(t, err, "属性 Name 的值为非空字符", result, false)
}

// isBlank的基准测试
func Benchmark_IsBlank(b *testing.B) {
	var value IsBlankEntity2
	for i := 0; i < b.N; i++ {
		value = IsBlankEntity2{Name: "zhou"}
		validate.Check(value, "name")
	}
}
