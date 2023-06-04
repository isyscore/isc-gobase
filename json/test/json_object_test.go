package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/json"
	"github.com/magiconair/properties/assert"
	"testing"
)

//{
//    "k1":12,
//    "k2":true,
//    "k3":{
//        "k31":32,
//        "k32":"str",
//        "k33":{
//            "k331":12
//        }
//    }
//}
var str = "{\n    \"k1\":12,\n    \"k2\":true,\n    \"k3\":{\n        \"k31\":32,\n        \"k32\":\"str\",\n        \"k33\":{\n            \"k331\":12\n        }\n    }\n}"

func TestLoad(t *testing.T) {
	jsonData := json.Object{}
	err := jsonData.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	assert.Equal(t, jsonData.Get("k1"), 12)
	assert.Equal(t, jsonData.Get("k2"), true)
	assert.Equal(t, jsonData.Get("k3.k31"), 32)
	assert.Equal(t, jsonData.Get("k3.k32"), "str")
	assert.Equal(t, jsonData.Get("k3.k33.k331"), 12)
}

func TestPut1(t *testing.T) {
	jsonObject := json.Object{}
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonObject.Put("k1", 13)
	jsonObject.Put("k2", false)
	jsonObject.Put("k3.k31", 33)
	jsonObject.Put("k3.k32", "str-change")
	jsonObject.Put("k3.k33.k331", 134)

	assert.Equal(t, jsonObject.Get("k1"), 13)
	assert.Equal(t, jsonObject.Get("k2"), false)
	assert.Equal(t, jsonObject.Get("k3.k31"), 33)
	assert.Equal(t, jsonObject.Get("k3.k32"), "str-change")
	assert.Equal(t, jsonObject.Get("k3.k33.k331"), 134)
}

func TestPut2(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\n    \"test\":12,\n    \"ok\":\"haode\",\n    \"k1\":{\n        \"k2\":true,\n        \"k21\":{\n            \"k3\":43,\n            \"k3array\":[\n                1,\n                2,\n                3,\n                4\n            ]\n        }\n    }\n}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonObject.Put("k1.k3", "value1")
	assert.Equal(t, jsonObject.Get("k1.k3"), "value1")
}

func TestPut3(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\n    \"test\":12,\n    \"ok\":\"haode\",\n    \"k1\":{\n        \"k2\":true,\n        \"k21\":{\n            \"k3\":43,\n            \"k3array\":[\n                1,\n                2,\n                3,\n                4\n            ]\n        }\n    }\n}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonObject.Put("k1.k4.k41", "value41")
	assert.Equal(t, jsonObject.Get("k1.k4.k41"), "value41")
}

func TestPut4(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\n    \"intValue\":12,\n    \"intValue8\":12,\n    \"intValue16\":12,\n    \"intValue32\":12,\n    \"intValue64\":12,\n    \"stringValue\":\"haode\",\n    \"boolValue\":false,\n    \"objectValue\":{\n        \"field1\":true,\n        \"field2Struct\":{\n            \"f21\":43,\n            \"k2array\":[\n                1,\n                2,\n                3,\n                4\n            ]\n        }\n    }\n}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	assert.Equal(t, jsonObject.GetInt("intValue"), 12)
	assert.Equal(t, jsonObject.GetInt8("intValue"), int8(12))
	assert.Equal(t, jsonObject.GetInt16("intValue"), int16(12))
	assert.Equal(t, jsonObject.GetInt32("intValue"), int32(12))
	assert.Equal(t, jsonObject.GetInt64("intValue"), int64(12))

	assert.Equal(t, jsonObject.GetString("stringValue"), "haode")
	assert.Equal(t, jsonObject.GetBool("boolValue"), false)
	assert.Equal(t, jsonObject.GetBool("objectValue.field1"), true)
	assert.Equal(t, jsonObject.GetInt("objectValue.field2Struct.f21"), 43)

	testEntity := TestEntity{}
	_ = jsonObject.GetObject("objectValue", &testEntity)

	fmt.Println(testEntity)
}

type TestEntity struct {
	Field1 bool
	Field2Struct TestEntity2
}

type TestEntity2 struct {
	F21 int
	K2array []int
}

// 普通数组
func TestGet5(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\"data\":{\"values\":[{\"name\":\"zhou\",\"age\":1},{\"name\":\"song\",\"age\":2}]}}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	assert.Equal(t, jsonObject.GetString("data.values[0].name"), "zhou")
	assert.Equal(t, jsonObject.GetString("data.values[0].age"), "1")
	assert.Equal(t, jsonObject.GetString("data.values[1].name"), "song")
	assert.Equal(t, jsonObject.GetString("data.values[1].age"), "2")
}

// 二维数组
func TestGet5_1(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\"data\":{\"values\":[[{\"name\":\"zhou\",\"age\":1},{\"name\":\"song\",\"age\":2}]]}}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	assert.Equal(t, jsonObject.GetString("data.values[0][0].name"), "zhou")
	assert.Equal(t, jsonObject.GetString("data.values[0][0].age"), "1")
	assert.Equal(t, jsonObject.GetString("data.values[0][1].name"), "song")
	assert.Equal(t, jsonObject.GetString("data.values[0][1].age"), "2")
}
