package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

// 测试基本表达式
type ValueConditionEntity1 struct {
	Data1 int `match:"condition=#current + #root.Data2 > 100"`
	Data2 int `match:"condition=#current < 20"`
	Data3 int `match:"condition=(++#current) >31"`
}

// 测试表达式
type ValueConditionEntity2 struct {
	Age   int `match:"condition=#root.Judge"`
	Judge bool
}

// 身份证号
func TestCondition(t *testing.T) {
	var value ValueConditionEntity1
	var result bool
	var err string

	// 测试 异常情况
	value = ValueConditionEntity1{Data1: 91, Data2: 10, Data3: 31}
	result,_ , err = validate.Check(value, "data1")
	True(t, result)

	// 测试 异常情况
	value = ValueConditionEntity1{Data1: 90, Data2: 10, Data3: 31}
	result,_ , err = validate.Check(value)
	Equal(t, err, "属性 Data1 的值 90 不符合条件 [#current + #root.Data2 > 100] ", result, false)

	// 测试 异常情况
	value = ValueConditionEntity1{Data1: 81, Data2: 20, Data3: 31}
	result,_ , err = validate.Check(value)
	Equal(t, err, "属性 Data2 的值 20 不符合条件 [#current < 20] ", result, false)

	// 测试 异常情况
	value = ValueConditionEntity1{Data1: 91, Data2: 10, Data3: 30}
	result,_ , err = validate.Check(value)
	Equal(t, err, "属性 Data3 的值 30 不符合条件 [(++#current) >31] ", result, false)
}

func TestCondition2(t *testing.T) {
	var value ValueConditionEntity2
	var result bool
	var err string

	// 测试 异常情况
	value = ValueConditionEntity2{Age: 12, Judge: true}
	result, _, err = validate.Check(value, "data1")
	True(t, result)

	// 测试 异常情况
	value = ValueConditionEntity2{Age: 12, Judge: false}
	result, _, err = validate.Check(value)
	Equal(t, err, "属性 Age 的值 12 不符合条件 [#root.Judge] ", result, false)
}
