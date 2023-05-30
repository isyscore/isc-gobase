package test

import (
	"encoding/json"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/test"
	"testing"
)

// 对以下的api进行测试
//
// mapToObject
// strToObject
// arrayToObject
// ToMap
// dataToObject：这个是总况
//
// objectToJson
// objectToData：这个是总的

// mapToObject
type ValueInnerEntity1 struct {
	Name string
	Age  int
}

type ValueInnerToMap1 struct {
	Name string `key:"name_test"`
	Age  int `key:"age_test"`
}

type ValueInnerToMap2 struct {
	Name string `key:"name_test"`
	Age  int `key:"age_test"`
	Address  int `key:"age_test ignore"`
}

func TestToMap1(t *testing.T) {
	var targetObj ValueInnerToMap1
	targetObj.Age = 12
	targetObj.Name = "test"

	inner1 := isc.ToMap(targetObj)
	test.Equal(t, isc.ToJsonString(inner1), "{\"age_test\":12,\"name_test\":\"test\"}")
}

func TestToMap2(t *testing.T) {
	var targetObj ValueInnerToMap1
	targetObj.Age = 12
	targetObj.Name = "test"

	inner1 := isc.ToMap(targetObj)
	test.Equal(t, isc.ToJsonString(inner1), "{\"age_test\":12,\"name_test\":\"test\"}")
}

func TestMapToObject1(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	var targetObj ValueInnerEntity1
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity1Extend struct {
	ValueInnerEntity1
	Address string
}

//func TestMapToObject1_extend(t *testing.T) {
//	inner1 := map[string]any{}
//	inner1["name"] = "inner_1"
//	inner1["age"] = 1
//	inner1["address"] = 1
//
//	var targetObj ValueInnerEntity1Extend
//	_ = isc.MapToObject(inner1, &targetObj)
//	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"Address\":1}", isc.ToJsonString(targetObj))
//}

type ValueInnerEntity2 struct {
	Name   string
	Age    int
	Inner1 ValueInnerEntity1
}

func TestMapToObject2(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	inner2 := map[string]any{}
	inner2["name"] = "inner_2"
	inner2["age"] = 2
	inner2["inner1"] = inner1

	var targetObj ValueInnerEntity2
	_ = isc.MapToObject(inner2, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_2\",\"Age\":2,\"Inner1\":{\"Name\":\"inner_1\",\"Age\":1}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity3 struct {
	Name   string
	Age    int
	Inner2 ValueInnerEntity2
}

func TestMapToObject3(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	inner2 := map[string]any{}
	inner2["name"] = "inner_2"
	inner2["age"] = 2
	inner2["inner1"] = inner1

	inner3 := map[string]any{}
	inner3["name"] = "inner_3"
	inner3["age"] = 3
	inner3["inner2"] = inner2

	var targetObj ValueInnerEntity3
	_ = isc.MapToObject(inner3, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_3\",\"Age\":3,\"Inner2\":{\"Name\":\"inner_2\",\"Age\":2,\"Inner1\":{\"Name\":\"inner_1\",\"Age\":1}}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity4 struct {
	Name    string
	Age     int
	DataMap map[string]string
}

func TestMapToObject4(t *testing.T) {
	kvMap := map[string]any{}
	kvMap["k1"] = "name1"
	kvMap["k2"] = "name2"

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity4
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":\"name1\",\"k2\":\"name2\"}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity5 struct {
	Name    string
	Age     int
	DataMap map[string]ValueInnerEntity1
}

func TestMapToObject5(t *testing.T) {
	v1 := map[string]any{}
	v1["name"] = "inner_1"
	v1["age"] = 1

	v2 := map[string]any{}
	v2["name"] = "inner_2"
	v2["age"] = 2

	kvMap := map[string]any{}
	kvMap["k1"] = v1
	kvMap["k2"] = v2

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity5
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":{\"Name\":\"inner_1\",\"Age\":1},\"k2\":{\"Name\":\"inner_2\",\"Age\":2}}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity6 struct {
	Name    string
	Age     int
	DataMap map[string][]int
}

func TestMapToObject6(t *testing.T) {
	var dataList []int
	dataList = append(dataList, 12)
	dataList = append(dataList, 13)

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity6
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[12,13],\"k2\":[12,13]}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity7 struct {
	Name    string
	Age     int
	DataMap map[string][]ValueInnerEntity1
}

func TestMapToObject7(t *testing.T) {
	var dataList []ValueInnerEntity1
	dataList = append(dataList, ValueInnerEntity1{Name: "name1", Age: 1})
	dataList = append(dataList, ValueInnerEntity1{Name: "name2", Age: 2})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity7
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity1Tem struct {
	Name    string
	Address string
}

type ValueInnerEntity8 struct {
	Name    string
	Age     int
	DataMap map[string][]ValueInnerEntity1Tem
}

func TestMapToObject8(t *testing.T) {
	var dataList []ValueInnerEntity1
	dataList = append(dataList, ValueInnerEntity1{Name: "name1", Age: 1})
	dataList = append(dataList, ValueInnerEntity1{Name: "name2", Age: 2})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity8
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Address\":\"\"},{\"Name\":\"name2\",\"Address\":\"\"}],\"k2\":[{\"Name\":\"name1\",\"Address\":\"\"},{\"Name\":\"name2\",\"Address\":\"\"}]}}", isc.ToJsonString(targetObj))
}

type ValueInnerEntity9Tem struct {
	Name string
	Age  string
}

type ValueInnerEntity9 struct {
	Name    string
	Age     int
	DataMap map[string][]ValueInnerEntity1
}

func TestMapToObject9(t *testing.T) {
	var dataList []ValueInnerEntity9Tem
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name1", Age: "1"})
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name2", Age: "2"})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity9
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

type ConfigValueTypeEnum int

const (
	YAML       ConfigValueTypeEnum = 0
	PROPERTIES ConfigValueTypeEnum = 1
	JSON       ConfigValueTypeEnum = 2
	STRING     ConfigValueTypeEnum = 3
)

type ValueInnerEntity10 struct {
	Name    string
	Age     ConfigValueTypeEnum
	DataMap map[string][]ValueInnerEntity1
}

func TestMapToObject10(t *testing.T) {
	var dataList []ValueInnerEntity9Tem
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name1", Age: "1"})
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name2", Age: "2"})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity9
	_ = isc.MapToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

func TestMapToObject11(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 12

	inner2 := map[string]any{}

	_ = isc.MapToObject(inner1, &inner2)
	test.Equal(t, isc.ToJsonString(inner2), "{\"age\":12,\"name\":\"inner_1\"}")
}

func TestMapToObject12(t *testing.T) {
	inner1 := map[string]string{}
	inner1["name"] = "inner_1"
	inner1["age"] = "12"

	inner2 := map[string]any{}

	_ = isc.MapToObject(inner1, &inner2)
	test.Equal(t, "{\"age\":\"12\",\"name\":\"inner_1\"}", isc.ToJsonString(inner2))
}

func TestMapToObject13(t *testing.T) {
	inner1 := map[string]any{}
	inner1["age"] = 12

	inner2 := map[string]int{}

	_ = isc.MapToObject(inner1, &inner2)
	test.Equal(t, "{\"age\":12}", isc.ToJsonString(inner2))
}

// dataToObject
func TestDataToObject1(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	var targetObj ValueInnerEntity1
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1}", isc.ToJsonString(targetObj))
}

func TestDataToObject2(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	inner2 := map[string]any{}
	inner2["name"] = "inner_2"
	inner2["age"] = 2
	inner2["inner1"] = inner1

	var targetObj ValueInnerEntity2
	_ = isc.DataToObject(inner2, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_2\",\"Age\":2,\"Inner1\":{\"Name\":\"inner_1\",\"Age\":1}}", isc.ToJsonString(targetObj))
}

func TestDataToObject3(t *testing.T) {
	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1

	inner2 := map[string]any{}
	inner2["name"] = "inner_2"
	inner2["age"] = 2
	inner2["inner1"] = inner1

	inner3 := map[string]any{}
	inner3["name"] = "inner_3"
	inner3["age"] = 3
	inner3["inner2"] = inner2

	var targetObj ValueInnerEntity3
	_ = isc.DataToObject(inner3, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_3\",\"Age\":3,\"Inner2\":{\"Name\":\"inner_2\",\"Age\":2,\"Inner1\":{\"Name\":\"inner_1\",\"Age\":1}}}", isc.ToJsonString(targetObj))
}

func TestDataToObject4(t *testing.T) {
	kvMap := map[string]any{}
	kvMap["k1"] = "name1"
	kvMap["k2"] = "name2"

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity4
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":\"name1\",\"k2\":\"name2\"}}", isc.ToJsonString(targetObj))
}

func TestDataToObject5(t *testing.T) {
	v1 := map[string]any{}
	v1["name"] = "inner_1"
	v1["age"] = 1

	v2 := map[string]any{}
	v2["name"] = "inner_2"
	v2["age"] = 2

	kvMap := map[string]any{}
	kvMap["k1"] = v1
	kvMap["k2"] = v2

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity5
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":{\"Name\":\"inner_1\",\"Age\":1},\"k2\":{\"Name\":\"inner_2\",\"Age\":2}}}", isc.ToJsonString(targetObj))
}

func TestDataToObject6(t *testing.T) {
	var dataList []int
	dataList = append(dataList, 12)
	dataList = append(dataList, 13)

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity6
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[12,13],\"k2\":[12,13]}}", isc.ToJsonString(targetObj))
}

func TestDataToObject7(t *testing.T) {
	var dataList []ValueInnerEntity1
	dataList = append(dataList, ValueInnerEntity1{Name: "name1", Age: 1})
	dataList = append(dataList, ValueInnerEntity1{Name: "name2", Age: 2})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity7
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

func TestDataToObject8(t *testing.T) {
	var dataList []ValueInnerEntity1
	dataList = append(dataList, ValueInnerEntity1{Name: "name1", Age: 1})
	dataList = append(dataList, ValueInnerEntity1{Name: "name2", Age: 2})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity8
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Address\":\"\"},{\"Name\":\"name2\",\"Address\":\"\"}],\"k2\":[{\"Name\":\"name1\",\"Address\":\"\"},{\"Name\":\"name2\",\"Address\":\"\"}]}}", isc.ToJsonString(targetObj))
}

func TestDataToObject9(t *testing.T) {
	var dataList []ValueInnerEntity9Tem
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name1", Age: "1"})
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name2", Age: "2"})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity9
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

func TestDataToObject10(t *testing.T) {
	var dataList []ValueInnerEntity9Tem
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name1", Age: "1"})
	dataList = append(dataList, ValueInnerEntity9Tem{Name: "name2", Age: "2"})

	kvMap := map[string]any{}
	kvMap["k1"] = dataList
	kvMap["k2"] = dataList

	inner1 := map[string]any{}
	inner1["name"] = "inner_1"
	inner1["age"] = 1
	inner1["dataMap"] = kvMap

	var targetObj ValueInnerEntity9
	_ = isc.DataToObject(inner1, &targetObj)
	test.Equal(t, "{\"Name\":\"inner_1\",\"Age\":1,\"DataMap\":{\"k1\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}],\"k2\":[{\"Name\":\"name1\",\"Age\":1},{\"Name\":\"name2\",\"Age\":2}]}}", isc.ToJsonString(targetObj))
}

// strToObject
func TestStrToObject1(t *testing.T) {
	var targetObj int
	_ = isc.StrToObject("123", &targetObj)
	test.Equal(t, targetObj, 123)
}

func TestStrToObject2(t *testing.T) {
	var targetObj string
	_ = isc.StrToObject("ok", &targetObj)
	test.Equal(t, targetObj, "ok")
}

func TestStrToObject3(t *testing.T) {
	var targetObj string
	_ = isc.StrToObject("{\"nihao\": \"haode\"}", &targetObj)
	test.Equal(t, targetObj, "{\"nihao\": \"haode\"}")
}

func TestStrToObject4(t *testing.T) {
	var targetObj ValueInnerEntity1
	_ = isc.StrToObject("{\"Age\": 12}", &targetObj)
	test.Equal(t, isc.ToJsonString(targetObj), "{\"Name\":\"\",\"Age\":12}")
}

func TestStrToObject5(t *testing.T) {
	var targetObj ValueInnerEntity1
	_ = isc.StrToObject("{\"age\": 12}", &targetObj)
	test.Equal(t, isc.ToJsonString(targetObj), "{\"Name\":\"\",\"Age\":12}")
}

func TestStrToObject6(t *testing.T) {
	targetObj := map[string]any{}
	_ = isc.StrToObject("{\"age\": 12}", &targetObj)
	test.Equal(t, isc.ToJsonString(targetObj), "{\"age\":12}")
}

func TestStrToObject7(t *testing.T) {
	var targetObj []ValueInnerEntity1
	_ = isc.StrToObject("[{\"Age\": 12},{\"Age\":14}]", &targetObj)
	test.Equal(t, isc.ObjectToJson(targetObj), "[{\"age\":12,\"name\":\"\"},{\"age\":14,\"name\":\"\"}]")
}

type ValueInnerEntityStr1 struct {
	//Name    string
	//Age     int
	DataMap any
}

func TestStrToObject8(t *testing.T) {
	str := "{\"dataMap\":{\"haha\":12,\"innerKey\":\"ok\"}}"

	var targetObj ValueInnerEntityStr1
	_ = isc.StrToObject(str, &targetObj)
	test.Equal(t, isc.ObjectToJson(targetObj), str)
}

// arrayToObject
func TestArrayToObject1(t *testing.T) {
	var dstValues []ValueInnerEntity1
	var targetObjs []ValueInnerEntity1
	targetObjs = append(targetObjs, ValueInnerEntity1{Name: "zhou", Age: 1})

	_ = isc.ArrayToObject(targetObjs, &dstValues)
	test.Equal(t, isc.ObjectToJson(dstValues), "[{\"age\":1,\"name\":\"zhou\"}]")
}

//todo 这个暂时还有点问题
//func TestArrayToObject2(t *testing.T) {
//	var dstArray []map[string]any
//	var srcArray []ValueInnerEntity1
//	srcArray = append(srcArray, ValueInnerEntity1{Name: "zhou", Age: 1})
//
//	_ = isc.ArrayToObject(srcArray, &dstArray)
//	Equal(t, isc.ObjectToJson(dstArray), "[{\"age\":1,\"name\":\"zhou\"}]")
//}

type ConfigItemFromCommonReq struct {
	Profile       string `match:"customize=ExistProfile" errMsg:"环境变量：#current 不存在或没有激活"`
	AppName       string
	ConfigItemKey string
}

func TestTtt(t *testing.T) {
	str := "{\"configItemKey\":null}"
	req := ConfigItemFromCommonReq{}
	_ = isc.StrToObject(str, &req)
	t.Log(req)
}

// objectToJson
type ValueObjectTest1 struct {
	AppName string
	Age     int
}

func TestObjectToJson1(t *testing.T) {
	entity := ValueObjectTest1{AppName: "zhou", Age: 12}
	test.Equal(t, isc.ObjectToJson(entity), "{\"age\":12,\"appName\":\"zhou\"}")
}

type ValueObjectTest2 struct {
	AppName string

	Age1 int
	Age2 int8
	Age3 int16
	Age4 int32
	Age5 int64

	UAge1 uint
	UAge2 uint8
	UAge3 uint16
	UAge4 uint32
	UAge5 uint64

	FAge1 float32
	FAge2 float64

	CAge1 complex64
	CAge2 complex128
}

func TestObjectToJson2(t *testing.T) {
	entity := ValueObjectTest2{
		AppName: "zhou",
		Age1:    12,
		Age2:    12,
		Age3:    12,
		Age4:    12,
		Age5:    12,
		UAge1:   12,
		UAge2:   12,
		UAge3:   12,
		UAge4:   12,
		UAge5:   12,
		FAge1:   12.1,
		FAge2:   12.2,
		CAge1:   3.2 + 12i,
		CAge2:   5.2 + 13i,
	}
	test.Equal(t, isc.ObjectToJson(entity), "{\"age1\":12,\"age2\":12,\"age3\":12,\"age4\":12,\"age5\":12,\"appName\":\"zhou\",\"cAge1\":\"(3.2+12i)\",\"cAge2\":\"(5.2+13i)\",\"fAge1\":12.1,\"fAge2\":12.2,\"uAge1\":12,\"uAge2\":12,\"uAge3\":12,\"uAge4\":12,\"uAge5\":12}")
}

type ValueObjectTest3 struct {
	AppName []string
	Age1    map[string]any
}

func TestObjectToJson3(t *testing.T) {
	var arrays []string
	arrays = append(arrays, "zhou")
	arrays = append(arrays, "wang")

	dataMap := map[string]any{}
	dataMap["a"] = 1
	dataMap["b"] = 2

	entity := ValueObjectTest3{
		AppName: arrays,
		Age1:    dataMap,
	}
	test.Equal(t, isc.ObjectToJson(entity), "{\"age1\":{\"a\":1,\"b\":2},\"appName\":[\"zhou\",\"wang\"]}")
}

type ValueObjectTest4 struct {
	AppName string
	Inner   ValueObjectTest3
}

func TestObjectToJson4(t *testing.T) {
	var arrays []string
	arrays = append(arrays, "zhou")
	arrays = append(arrays, "wang")

	dataMap := map[string]any{}
	dataMap["a"] = 1
	dataMap["b"] = 2

	entity3 := ValueObjectTest3{
		AppName: arrays,
		Age1:    dataMap,
	}

	var entity4 ValueObjectTest4
	entity4.Inner = entity3
	entity4.AppName = "zhou"
	test.Equal(t, isc.ObjectToJson(entity4), "{\"appName\":\"zhou\",\"inner\":{\"age1\":{\"a\":1,\"b\":2},\"appName\":[\"zhou\",\"wang\"]}}")
}

func TestObjectToJson5(t *testing.T) {
	var arrays []string
	arrays = append(arrays, "zhou")
	arrays = append(arrays, "wang")

	dataMap := map[string]any{}
	dataMap["A"] = 1
	dataMap["B"] = 2

	act := ValueObjectTest3{
		AppName: arrays,
		Age1:    dataMap,
	}
	test.Equal(t, isc.ObjectToJson(act), "{\"age1\":{\"a\":1,\"b\":2},\"appName\":[\"zhou\",\"wang\"]}")
}

func TestObjectToJson6(t *testing.T) {
	expect := "[1,2]"
	var act []int
	act = append(act, 1)
	act = append(act, 2)
	test.Equal(t, isc.ObjectToJson(act), expect)
}

func TestObjectToJson7(t *testing.T) {
	var act []ValueInnerEntity1
	act = append(act, ValueInnerEntity1{Name: "zhou1", Age: 1})
	act = append(act, ValueInnerEntity1{Name: "zhou2", Age: 2})
	expect := "[{\"age\":1,\"name\":\"zhou1\"},{\"age\":2,\"name\":\"zhou2\"}]"
	test.Equal(t, isc.ObjectToJson(act), expect)
}

func TestObjectToJson8(t *testing.T) {
	var act = []map[string]any{}

	map1 := map[string]any{}
	map1["name"] = "zhou1"
	map1["age"] = 1

	map2 := map[string]any{}
	map2["name"] = "zhou2"
	map2["age"] = 2

	act = append(act, map1)
	act = append(act, map2)
	test.Equal(t, isc.ObjectToJson(act), "[{\"age\":1,\"name\":\"zhou1\"},{\"age\":2,\"name\":\"zhou2\"}]")
}

type PageRsp struct {

	// 分页数据
	Records []any
}

func TestObjectToJson9(t *testing.T) {
	rel := "{\"Records\":[{\"Id\":121,\"AppName\":\"asdf\",\"AppDesc\":\"fffds\",\"ActiveStatus\":1,\"CreateTime\":\"2021-12-20 14:05:10 +0800 CST\",\"UpdateTime\":\"2021-12-21 14:19:13 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":0,\"Version\":0},{\"Id\":117,\"AppName\":\"isc-apaas-service\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-12-06 10:32:11 +0800 CST\",\"UpdateTime\":\"2021-12-06 10:32:11 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":10,\"Version\":0},{\"Id\":116,\"AppName\":\"isc-config-sample-3\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-11-05 19:35:07 +0800 CST\",\"UpdateTime\":\"2021-11-05 19:35:07 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":4,\"Version\":0},{\"Id\":115,\"AppName\":\"isc-config-sample-2\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-11-05 19:27:55 +0800 CST\",\"UpdateTime\":\"2021-11-05 19:27:55 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":4,\"Version\":0},{\"Id\":113,\"AppName\":\"isc-config-sample1\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-11-05 19:21:18 +0800 CST\",\"UpdateTime\":\"2021-11-05 19:21:18 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":0,\"Version\":0},{\"Id\":112,\"AppName\":\"app-demo-xxx\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-28 14:29:46 +0800 CST\",\"UpdateTime\":\"2021-09-28 14:29:46 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":3,\"Version\":25},{\"Id\":84,\"AppName\":\"isc-config-sample3\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-22 13:46:36 +0800 CST\",\"UpdateTime\":\"2021-09-22 13:46:36 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":2,\"Version\":23},{\"Id\":83,\"AppName\":\"isc-config-sample-local\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-22 11:22:12 +0800 CST\",\"UpdateTime\":\"2021-09-22 11:22:12 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":29,\"Version\":22},{\"Id\":82,\"AppName\":\"isc-monitoring-service2\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-18 11:47:50 +0800 CST\",\"UpdateTime\":\"2021-09-18 11:47:50 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":15,\"Version\":0},{\"Id\":81,\"AppName\":\"isc-monitoring-service1\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-11 18:02:25 +0800 CST\",\"UpdateTime\":\"2021-09-11 18:02:25 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":15,\"Version\":21},{\"Id\":80,\"AppName\":\"lamp-demo-a\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-10 17:21:40 +0800 CST\",\"UpdateTime\":\"2021-09-10 17:28:18 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":1,\"Version\":20},{\"Id\":79,\"AppName\":\"pivotdemoa\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-10 17:17:33 +0800 CST\",\"UpdateTime\":\"2021-09-10 17:17:33 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":1,\"Version\":19},{\"Id\":78,\"AppName\":\"isc-config-sample2\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-09 16:16:55 +0800 CST\",\"UpdateTime\":\"2021-09-09 16:16:55 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":25,\"Version\":18},{\"Id\":77,\"AppName\":\"isc-config-sample-client\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-09-09 13:40:10 +0800 CST\",\"UpdateTime\":\"2021-09-09 13:40:10 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":25,\"Version\":17},{\"Id\":76,\"AppName\":\"isc-pivot-client\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-08-30 17:53:10 +0800 CST\",\"UpdateTime\":\"2021-08-31 10:08:53 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":2,\"Version\":16},{\"Id\":74,\"AppName\":\"isc-rpc-3-os0\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-08-30 17:52:17 +0800 CST\",\"UpdateTime\":\"2021-08-30 17:52:17 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":25,\"Version\":14},{\"Id\":73,\"AppName\":\"isc-rpc-os0\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-07-06 15:23:25 +0800 CST\",\"UpdateTime\":\"2021-07-06 15:23:25 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":25,\"Version\":13},{\"Id\":71,\"AppName\":\"isc-common-service-test\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-06-21 21:50:24 +0800 CST\",\"UpdateTime\":\"2021-06-21 21:50:24 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":63,\"Version\":11},{\"Id\":70,\"AppName\":\"isc-config-sample\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-06-21 15:33:39 +0800 CST\",\"UpdateTime\":\"2021-06-21 15:33:39 +0800 CST\",\"CreateUser\":\"\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":30,\"Version\":10},{\"Id\":68,\"AppName\":\"isc-route-service1\",\"AppDesc\":\"\",\"ActiveStatus\":1,\"CreateTime\":\"2021-06-09 16:00:38 +0800 CST\",\"UpdateTime\":\"2021-06-09 16:00:38 +0800 CST\",\"CreateUser\":\"admin\",\"UpdateUser\":\"\",\"MachineNum\":0,\"ConfigNum\":1,\"Version\":8}]}\n"

	rp := PageRsp{}
	_ = isc.DataToObject(rel, &rp)

	t.Log(isc.ToJsonString(rp))

	result := isc.ObjectToData(rp)
	t.Log(isc.ToJsonString(result))

	//Equal(t, isc.ObjectToJson(act), "[{\"age\":1,\"name\":\"zhou1\"},{\"age\":2,\"name\":\"zhou2\"}]")
}

type JsonEntity10Inner struct {
	Name string
}

type JsonEntity10 struct {
	Data *JsonEntity10Inner
}

func TestObjectToJson10(t *testing.T) {
	inner := JsonEntity10Inner{
		Name: "ok",
	}
	entity := JsonEntity10{
		Data: &inner,
	}

	test.Equal(t, isc.ObjectToJson(entity), "{\"data\":{\"name\":\"ok\"}}")
}

func TestObjectToJson11(t *testing.T) {
	inner := JsonEntity10Inner{
		Name: "ok",
	}

	map1 := map[string]any{}
	map1["data"] = &inner

	test.Equal(t, isc.ObjectToJson(map1), "{\"data\":{\"name\":\"ok\"}}")
}

// objectToMap

// objectToArray

// objectToData
func TestObjectToData1(t *testing.T) {
	test.Equal(t, isc.ObjectToData(1), 1)
}

func TestObjectToData2(t *testing.T) {
	test.Equal(t, isc.ObjectToData("12"), "12")
}

func TestObjectToData3(t *testing.T) {
	test.Equal(t, isc.ObjectToData("ab"), "ab")
}

func TestObjectToData4(t *testing.T) {
	test.Equal(t, isc.ObjectToData(12.4), 12.4)
}

func TestObjectToData5(t *testing.T) {
	src := ValueInnerEntity1{Name: "zhou", Age: 12}
	dst := map[string]any{}
	dst["name"] = "zhou"
	dst["age"] = 12
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData6(t *testing.T) {
	src := map[string]any{}
	src["name"] = "zhou"
	src["age"] = 12

	dst := ValueInnerEntity1{Name: "zhou", Age: 12}
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData7(t *testing.T) {
	src := map[string]any{}
	src["name"] = "zhou"
	src["age"] = 12

	dst := map[string]any{}
	dst["name"] = "zhou"
	dst["age"] = 12
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData8(t *testing.T) {
	src := ValueInnerEntity1{Name: "zhou", Age: 12}
	dst := ValueInnerEntity1{Name: "zhou", Age: 12}
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

//type ValueInnerEntity1Json struct {
//	Age  int
//	Address string
//}

//func TestObjectToData9(t *testing.T) {
//	src := ValueInnerEntity1{Name: "zhou", Age: 12}
//	dst := ValueInnerEntity1Json{Age: 12}
//	Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
//}

func TestObjectToData10(t *testing.T) {
	src := ValueInnerEntity1{Name: "zhou", Age: 12}
	dst := map[string]any{}
	dst["name"] = "zhou"
	dst["age"] = 12
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData11(t *testing.T) {
	var src []ValueInnerEntity1
	var dst []ValueInnerEntity1
	src = append(src, ValueInnerEntity1{Name: "zhou", Age: 12})
	dst = append(dst, ValueInnerEntity1{Name: "zhou", Age: 12})
	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData12(t *testing.T) {
	var src []ValueInnerEntity1
	var dst []map[string]any
	src = append(src, ValueInnerEntity1{Name: "zhou", Age: 12})

	map1 := map[string]any{}
	map1["name"] = "zhou"
	map1["age"] = 12
	dst = append(dst, map1)

	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func TestObjectToData13(t *testing.T) {
	var dst []ValueInnerEntity1
	var src []map[string]any
	dst = append(dst, ValueInnerEntity1{Name: "zhou", Age: 12})

	map1 := map[string]any{}
	map1["name"] = "zhou"
	map1["age"] = 12
	src = append(src, map1)

	test.Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
}

func BenchmarkSprintfPress(b *testing.B) {
	//func TestSimple(t *testing.T) {
	var jsonStr = "{\"age\":12,\"name\":\"zhou\"}"

	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var entity ValueInnerEntity1
		_ = isc.DataToObject(jsonStr, &entity)

		resultMap := make(map[string]any)
		_ = json.Unmarshal([]byte(jsonStr), &resultMap)
	}
}

//type ValueInnerEntityPtr struct {
//	Ptr *ValueInnerEntity1
//}

//func TestObjectToData14(t *testing.T) {
//	entity := ValueInnerEntity1{
//		Name: "zhou",
//		Age:  12,
//	}
//
//	act := ValueInnerEntityPtr{}
//
//	map1 := map[string]any{}
//	map1["ptr"] = &entity
//
//	isc.MapToObject(map1, &act)
//	fmt.Println(act.Ptr.Name)
//	fmt.Println(act.Ptr.Age)
//
//	//Equal(t, isc.ObjectToJson(isc.ObjectToData(src)), isc.ObjectToJson(dst))
//}
