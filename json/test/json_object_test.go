package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/json"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\n    \"test\":12,\n    \"ok\":\"haode\",\n    \"k1\":{\n        \"k2\":true,\n        \"k21\":{\n            \"k3\":43,\n            \"k3array\":[\n                1,\n                2,\n                3,\n                4\n            ]\n        }\n    }\n}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	assert.Equal(t, 12, jsonObject.Get("test"))
	assert.Equal(t, "haode", jsonObject.Get("ok"))
	assert.Equal(t, true, jsonObject.Get("k1.k2"))
	assert.Equal(t, 43, jsonObject.Get("k1.k21.k3"))
}

func TestPut1(t *testing.T) {
	jsonObject := json.Object{}
	str := "{\n    \"test\":12,\n    \"ok\":\"haode\",\n    \"k1\":{\n        \"k2\":true,\n        \"k21\":{\n            \"k3\":43,\n            \"k3array\":[\n                1,\n                2,\n                3,\n                4\n            ]\n        }\n    }\n}"
	err := jsonObject.Load(str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonObject.Put("test", 23)
	jsonObject.Put("ok", "ok-change")
	jsonObject.Put("k1.k2", false)

	assert.Equal(t, jsonObject.Get("test"), 23)
	assert.Equal(t, jsonObject.Get("ok"), "ok-change")
	assert.Equal(t, jsonObject.Get("k1.k2"), false)
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

	testEntity := TestEntity{}
	jsonObject.GetObject("objectValue", &testEntity)

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
