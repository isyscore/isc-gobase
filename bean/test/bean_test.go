package test

import (
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAddBean(t *testing.T) {
	tt := TestEntity{Name: "hello", Age: 12}
	bean.AddBean("test", &tt)

	t1 := bean.GetBean("test")
	t2 := t1.(*TestEntity)
	assert.Equal(t, t2.Name, tt.Name)

	bean.Clean()
}

func TestGetFieldShow(t *testing.T) {
	tt := TestEntity{Name: "value"}
	// 添加注册
	bean.AddBean("test", &tt)

	// 获取值
	actValue := bean.GetField("test", "Name")
	assert.Equal(t, actValue, "value")

	// 修改值
	bean.SetField("test", "Name", "value-change")

	// 查看
	actValue = bean.GetField("test", "Name")
	assert.Equal(t, actValue, "value-change")

	bean.Clean()
}

func TestGetField(t *testing.T) {
	tt := TestEntity{Name: "hello", Age: 12}
	bean.AddBean("test", &tt)

	fv := bean.GetField("test", "Name")
	assert.Equal(t, fv, "hello")

	bean.Clean()
}

func TestExist(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	assert.Equal(t, bean.ExistBean("test"), true)

	bean.Clean()
}

func TestCallFun(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	fv := bean.CallFun("test", "Fun", map[string]any{})
	assert.Equal(t, isc.ToString(fv[0]), "ok")

	bean.Clean()
}

func TestCallFunPtr(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	fv := bean.CallFun("test", "Fun", map[string]any{})
	assert.Equal(t, isc.ToString(fv[0]), "ok")

	bean.Clean()
}

func TestCallFunUpper(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	fv := bean.CallFun("test", "Fun", map[string]any{})
	assert.Equal(t, isc.ToString(fv[0]), "ok")

	bean.Clean()
}

func TestCallFun1(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)

	parameterMap := map[string]any{}
	parameterMap["p1"] = "name"

	fv := bean.CallFun("test", "Fun1", parameterMap)
	assert.Equal(t, isc.ToString(fv[0]), "name")

	bean.Clean()
}

func TestCallFun2(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)

	parameterMap := map[string]any{}
	parameterMap["p1"] = 12

	fv := bean.CallFun("test", "Fun2", parameterMap)
	assert.Equal(t, isc.ToInt(fv[0]), 12)

	bean.Clean()
}

func TestCallFun3(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)

	parameterMap := map[string]any{}
	parameterMap["p1"] = "name"
	parameterMap["p2"] = 12

	fv := bean.CallFun("test", "Fun3", parameterMap)
	assert.Equal(t, isc.ToInt(fv[0]), 12)

	bean.Clean()
}

//// error: 暂时参数中的json到结构体的转换
//func TestCallFun4(t *testing.T) {
//	tt := TestEntity{}
//	inner := TestInnerEntity{Address: "杭州"}
//
//	bean.AddBean("test", &tt)
//
//	parameterMap := map[string]any{}
//	parameterMap["p1"] = "name"
//	parameterMap["p2"] = 12
//	parameterMap["p3"] = inner
//
//	fv := bean.CallFun("test", "Fun4", parameterMap)
//	assert.Equal(t, isc.ToString(fv[0]), "杭州")
//
//	bean.Clean()
//}

//// error: 暂时参数中的json到结构体的转换
//func TestCallFun4_1(t *testing.T) {
//	tt := TestEntity{}
//
//	bean.AddBean("test", &tt)
//
//	parameterMap := map[string]any{}
//	parameterMap["p1"] = "name"
//	parameterMap["p2"] = 12
//	parameterMap["p3"] = "{\"Address\": \"hangzhou\"}"
//
//	fv := bean.CallFun("test", "Fun4", parameterMap)
//	assert.Equal(t, isc.ToString(fv[0]), "杭州")
//
//	bean.Clean()
//}
//
//func TestCallFun4Ptr(t *testing.T) {
//	tt := TestEntity{}
//	inner := TestInnerEntity{Address: "杭州"}
//
//	bean.AddBean("test", &tt)
//
//	parameterMap := map[string]any{}
//	parameterMap["p1"] = "name"
//	parameterMap["p2"] = 12
//	parameterMap["p3"] = &inner
//
//	fv := bean.CallFun("test", "Fun4Ptr", parameterMap)
//	assert.Equal(t, isc.ToString(fv[0]), "杭州")
//
//	bean.Clean()
//}

func TestCallFun5(t *testing.T) {
	tt := TestEntity{Age: 12}

	bean.AddBean("test", &tt)

	fv := bean.CallFun("test", "Fun5", map[string]any{})
	assert.Equal(t, isc.ToInt(fv[0]), 12)

	bean.Clean()
}

func TestCallPtrFun(t *testing.T) {
	tt := TestEntity{Age: 12}

	bean.AddBean("test", &tt)

	fv := bean.CallFun("test", "PtrFun", map[string]any{})
	assert.Equal(t, isc.ToString(fv[0]), "ok")

	bean.Clean()
}

func TestCallPtrFun1(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)

	parameterMap := map[string]any{}
	parameterMap["p1"] = "name"

	fv := bean.CallFun("test", "PtrFun1", parameterMap)
	assert.Equal(t, isc.ToString(fv[0]), "name")

	bean.Clean()
}

func TestSetField(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	bean.SetField("test", "Name", "hello")

	assert.Equal(t, tt.Name, "hello")

	bean.Clean()
}

func TestSetField1(t *testing.T) {
	tt := TestEntity{Age: 12}

	bean.AddBean("test", &tt)
	parameterMap := map[string]any{}
	parameterMap["p1"] = 32

	assert.Equal(t, tt.Age, 12)
	bean.CallFun("test", "ChangeField", parameterMap)
	assert.Equal(t, tt.Age, 32)

	bean.Clean()
}

func TestSetFieldPtr(t *testing.T) {
	tt := TestEntity{}

	bean.AddBean("test", &tt)
	bean.SetField("test", "Name", "hello")

	assert.Equal(t, tt.Name, "hello")

	bean.Clean()
}

func TestGetBeans(t *testing.T) {
	tt1 := TestEntity{}
	tt2 := TestEntity{}
	tt3 := TestEntity{}

	bean.AddBean("t1-name", &tt1)
	bean.AddBean("t2-name2", &tt2)
	bean.AddBean("t3-change", &tt3)

	datas := bean.GetBeanNames("name")
	assert.Equal(t, len(datas), 2)

	datas = bean.GetBeanNames("")
	assert.Equal(t, len(datas), 3)

	bean.Clean()
}

func TestServerBean(t *testing.T) {
	tt1 := TestEntity{Name: "t1", Age: 12}
	tt2 := TestEntity{Name: "t2", Age: 13}
	tt3 := TestEntity{Name: "t2", Age: 13}

	bean.AddBean("t1-name", &tt1)
	bean.AddBean("t2-name2", &tt2)
	bean.AddBean("t3-change", &tt3)

	server.Run()
}

type ValueInnerEntity struct {
	Name string
	Age  int
}

type TestEntity struct {
	Name string
	Age  int
}

type TestInnerEntity struct {
	Name    string
	Address string
}

func (tt TestEntity) Fun() string {
	return "ok"
}

func (tt TestEntity) Fun1(name string) string {
	return name
}

func (tt TestEntity) Fun2(age int) int {
	return age
}

func (tt TestEntity) Fun3(name string, age int) int {
	return age
}

func (tt TestEntity) Fun4(name string, age int, inner TestInnerEntity) string {
	return inner.Address
}

func (tt TestEntity) Fun4Ptr(name string, age int, inner *TestInnerEntity) string {
	return inner.Address
}

func (tt TestEntity) Fun5() int {
	return tt.Age
}

func (tt TestEntity) Fun6(age int) {
	tt.Age = age
}

func (tt *TestEntity) PtrFun() string {
	return "ok"
}

func (tt *TestEntity) PtrFun1(name string) string {
	return name
}

func (tt *TestEntity) PtrFun2(age int) int {
	return age
}

func (tt *TestEntity) ChangeField(age int) {
	tt.Age = age
}
