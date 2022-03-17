package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type ErrMsgEntity1 struct {
	Name string `match:"value=zhou" errMsg:"对应的值不合法"`
	Age  int
}

type ErrMsgEntity2 struct {
	Name string `match:"value=zhou"`
	Age  int    `match:"condition=#current > 10" errMsg:"当前的值不合法，应该大于10，当前值为#current，对应的名字为#root.Name"`
}

type ErrMsgEntity3 struct {
	Name string `match:"value=zhou" errMsg:"当前值不合法，只可为zhou，当前的值为#current，年龄为#root.Age"`
	Age  int    `match:"condition=#current > 10" errMsg:"当前的值不合法，应该大于10，当前值为#current，对应的名字为#root.Name"`
}

func TestErrMsg1(t *testing.T) {
	var value ErrMsgEntity1
	var result bool
	var err string

	// 测试 正常情况
	value = ErrMsgEntity1{Name: "zhou"}
	result, _ = validate.Check(value, "name")
	True(t, result)

	// 测试 正常情况
	value = ErrMsgEntity1{Name: "宋江"}
	result, err = validate.Check(value, "name")
	Equal(t, err, "对应的值不合法", result, false)
}

func TestErrMsg2(t *testing.T) {
	var value ErrMsgEntity2
	var result bool
	var err string

	// 测试 正常情况
	//value = ErrMsgEntity2{Name: "zhou", Age: 12}
	//result, _ = validate.Check(value)
	//assert.True(t, result)

	// 测试 正常情况
	value = ErrMsgEntity2{Name: "zhou", Age: 2}
	result, err = validate.Check(value)
	Equal(t, err, "当前的值不合法，应该大于10，当前值为2，对应的名字为zhou", result, false)
}

func TestErrMsg3(t *testing.T) {
	var value ErrMsgEntity3
	var result bool
	var err string

	// 测试 正常情况
	value = ErrMsgEntity3{Name: "zhou", Age: 12}
	result, _ = validate.Check(value)
	True(t, result)

	// 测试 正常情况
	value = ErrMsgEntity3{Name: "zhou", Age: 2}
	result, err = validate.Check(value)
	Equal(t, err, "当前的值不合法，应该大于10，当前值为2，对应的名字为zhou", result, false)

	// 测试 正常情况
	value = ErrMsgEntity3{Name: "宋江", Age: 12}
	result, err = validate.Check(value, "name")
	Equal(t, err, "当前值不合法，只可为zhou，当前的值为宋江，年龄为12", result, false)

	// 测试 正常情况
	value = ErrMsgEntity3{Name: "宋江", Age: 3}
	result, err = validate.Check(value, "age")
	Equal(t, err, "当前的值不合法，应该大于10，当前值为3，对应的名字为宋江", result, false)
}
