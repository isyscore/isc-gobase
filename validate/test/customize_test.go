package main

import (
	"encoding/json"
	"github.com/isyscore/isc-gobase/validate"
	"github.com/isyscore/isc-gobase/validate/test/fun"
	"testing"
)

func TestCustomize1(t *testing.T) {
	var value fun.CustomizeEntity1
	var result bool
	var err string

	// 测试 正常情况
	value = fun.CustomizeEntity1{Name: "zhou"}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity1{Name: "宋江"}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity1{Name: "陈真"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Name 的值 陈真 没命中只允许条件回调 [judge1Name] ", result, false)
}

func TestCustomize2(t *testing.T) {
	var value fun.CustomizeEntity2
	var result bool
	var err string

	// 测试 正常情况
	value = fun.CustomizeEntity2{Name: "zhou"}
	result, err = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity2{Name: "宋江"}
	result, err = validate.Check(value, "name")
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity2{Name: "陈真"}
	result, err = validate.Check(value)
	Equal(t, err, "没有命中可用的值'zhou'和'宋江'", result, false)
}

func TestCustomize3(t *testing.T) {
	var value fun.CustomizeEntity3
	var result bool
	var err string

	// 测试 正常情况
	value = fun.CustomizeEntity3{Name: "zhou", Age: 20}
	result, err = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity3{Name: "宋江", Age: 20}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity3{Name: "陈真"}
	result, err = validate.Check(value)
	Equal(t, err, "没有命中可用的值'zhou'和'宋江'", result, false)

	// 测试 正常情况
	value = fun.CustomizeEntity3{Name: "zhou", Age: 13}
	result, _ = validate.Check(value)
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity3{Name: "zhou", Age: 10}
	result, err = validate.Check(value)
	Equal(t, err, "用户[zhou]没有满足年龄age > 12，当前年龄为：10", result, false)
}

func TestCustomize4(t *testing.T) {
	var value fun.CustomizeEntity4
	var result bool
	var err string

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 20}
	result, err = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "宋江", Age: 20}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity4{Name: "陈真"}
	result, err = validate.Check(value)
	Equal(t, err, "没有命中可用的值'zhou'和'宋江'", result, false)

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 13}
	result, _ = validate.Check(value)
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 10}
	result, err = validate.Check(value)
	Equal(t, err, "用户[zhou]没有满足年龄age > 12，当前年龄为：10", result, false)
}

func TestCustomize5(t *testing.T) {
	var value fun.CustomizeEntity4
	var result bool
	var err string

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 20}
	result, err = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "宋江", Age: 20}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity4{Name: "陈真"}
	result, err = validate.Check(value)
	Equal(t, err, "没有命中可用的值'zhou'和'宋江'", result, false)

	// 测试 正常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 13}
	result, _ = validate.Check(value)
	True(t, result)

	// 测试 异常情况
	value = fun.CustomizeEntity4{Name: "zhou", Age: 10}
	result, err = validate.Check(value)
	Equal(t, err, "用户[zhou]没有满足年龄age > 12，当前年龄为：10", result, false)
}

func TestCustomize5_1(t *testing.T) {
	var value fun.CustomizeEntity5
	var result bool

	// 测试 正常情况
	value = fun.CustomizeEntity5{Name: "zhou", Age: 20}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity5{Name: "宋江", Age: 20}
	result, _ = validate.Check(value, "name")
	True(t, result)
}

func TestCustomize6(t *testing.T) {
	var value fun.CustomizeEntity6
	var result bool
	var pMap map[string]interface{}

	// 测试 正常情况
	value = fun.CustomizeEntity6{}
	pMap = map[string]interface{}{
		"name": "zhou",
		"age":  20,
	}
	result, _ = validate.CheckWithParameter(pMap, value, "name1")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity6{Name2: "zhou"}
	pMap = map[string]interface{}{
		"age": 20,
	}
	result, _ = validate.CheckWithParameter(pMap, value, "name2")
	True(t, result)

	// 测试 正常情况
	value = fun.CustomizeEntity6{Name3: "zhou"}
	pMap = map[string]interface{}{
		"age": 20,
	}
	result, _ = validate.CheckWithParameter(pMap, value, "name3")
	True(t, result)
}

func TestCustomize7(t *testing.T) {
	var value fun.CustomizeEntity7
	var result bool

	// 测试 正常情况
	value = fun.CustomizeEntity7{Name: nil}
	result, _ = validate.Check(value, "name")
	False(t, result)

	// 测试 正常情况
	//name := "df"
	//value = fun.CustomizeEntity7{Name: &name}
	//result, _ = validate.Check(value, "name")
	//True(t, result)
}

func TestCustomize7_1(t *testing.T) {
	var value fun.CustomizeEntity7
	var value1 fun.CustomizeEntity7
	var result bool

	// 测试 正常情况
	value = fun.CustomizeEntity7{Flag: nil}
	result, msg := validate.Check(value, "flag")
	FalseMsg(t, result, msg)

	// 测试 正常情况
	flag := true
	value = fun.CustomizeEntity7{Flag: &flag}
	result, _ = validate.Check(value, "flag")
	True(t, result)

	str := "{\"name\":\"xxx\", \"age\":12}"
	_ = json.Unmarshal([]byte(str), &value1)
	result, msg = validate.Check(value1, "flag")
	FalseMsg(t, result, msg)

	value = fun.CustomizeEntity7{Flag2: nil}
	result, msg = validate.Check(value, "flag2")
	FalseMsg(t, result, msg)

	flag = true
	value = fun.CustomizeEntity7{Flag2: &flag}
	result, _ = validate.Check(value, "flag2")
	True(t, result)
}
