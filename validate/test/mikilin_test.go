package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type MikilinBaseEntity struct {
	Name string `match:"value={zhou, 宋江} isBlank=true"`
	Age  int    `match:"value={12, 13}"`
}

type MikilinBaseEntity2 struct {
	Name string `match:"value={zhou, 宋江} isBlank"`
	Age  int    `match:"value={12, 13}"`
}

func TestMkBase1(t *testing.T) {
	var value MikilinBaseEntity
	var result bool
	var err string

	//测试 正常情况
	value = MikilinBaseEntity{Age: 12}
	result, _, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 正常情况
	value = MikilinBaseEntity{Age: 13, Name: "zhou"}
	result, _, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = MikilinBaseEntity{Age: 13, Name: "陈真"}
	result, _, err = validate.Check(value)
	Equal(t, err, "[\"属性 Name 的值 陈真 不在只可用列表 [zhou 宋江] 中\",\"属性 Name 的值为非空字符\"]", result, false)
}

func TestMkBase2(t *testing.T) {
	var value MikilinBaseEntity2
	var result bool
	var err string

	//测试 正常情况
	value = MikilinBaseEntity2{Age: 12}
	result, _, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 正常情况
	value = MikilinBaseEntity2{Age: 13, Name: "zhou"}
	result, _, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = MikilinBaseEntity2{Age: 13, Name: "陈真"}
	result, _, err = validate.Check(value)
	Equal(t, "[\"属性 Name 的值 陈真 不在只可用列表 [zhou 宋江] 中\",\"属性 Name 的值为非空字符\"]", err, result, false)
}

// 压测进行基准测试
func BenchmarkMkBase3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validate.Check(MikilinBaseEntity{Age: 12})
	}
}
