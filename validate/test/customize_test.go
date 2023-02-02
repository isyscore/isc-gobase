package main

import (
	"encoding/json"
	"fmt"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/validate"
	"github.com/isyscore/isc-gobase/validate/test/fun"
	"math"
	"math/rand"
	"testing"
	t0 "time"
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

func TestFun(t *testing.T) {
	rand.Seed(t0.Now().UnixNano())
	////随机生成100以内的正整数
	//
	//a := 1
	//b := 3
	//
	//// [-1.0, 1.0)
	//fmt.Println((isc.ToFloat64(a))/(isc.ToFloat64(b)*1.0))
	//
	//Solos := []*Solo{}
	//
	//Solos = append(Solos, &Solo{score: 12.0})
	//Solos = append(Solos, &Solo{score: 1.0})
	//Solos = append(Solos, &Solo{score: 132.2})
	//Solos = append(Solos, &Solo{score: 54.2})
	//Solos = append(Solos, &Solo{score: 32.2})
	//
	//group := Group{solos: Solos}
	//
	//sort.Sort(group)
	//
	//for _, solo := range group.solos {
	//	fmt.Println(solo.score)
	//}
	//
	//rand.Shuffle(len(group.solos), func(i, j int) {
	//	group.solos[i], group.solos[j] = group.solos[j], group.solos[i]
	//})
	//
	//fmt.Println("======")
	//for _, solo := range group.solos {
	//	fmt.Println(solo.score)
	//}
	//
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//fmt.Println(rand.Intn(2))
	//// 0,1
	//fmt.Println(rand.Intn(2))


	//datas := []int{1,2, 3, 4, 5}
	//fmt.Println(datas[0:2])
	//fmt.Println(datas[2:4])
	//fmt.Println(datas[:len(datas)])

	ratio := 0.3
	num := 12

	fmt.Println(isc.ToInt(math.Ceil(isc.ToFloat64(num) * ratio)))
}

// 个体
type Solo struct {
	// 评分
	score float64
}

// 种群
type Group struct {
	// 参数
	solos []*Solo
}

func (group Group) Len() int {
	return len(group.solos)
}

func (group Group) Less(i, j int) bool {
	return group.solos[i].score > group.solos[j].score
}

func (group Group) Swap(i, j int) {
	group.solos[i], group.solos[j] = group.solos[j], group.solos[i]
}
