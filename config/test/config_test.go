package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/magiconair/properties/assert"
	"testing"
)

//key1:
//  key2:
//    intdata: 12
//    strdata: "data"
//    booldata: true
//    int64data: 12
//    floatdata: 12.3
//    objdata:
//      field1: 1
//      field2: "value2"
//    arraydata:
//      - 1
//      - 2
//      - 3
//    arrayobjdata:
//      - field1: 1
//        field2: "name1"
//      - field1: 2
//        field2: "name2"
func TestLoadConfig(t *testing.T) {
	config.LoadConfig()

	assert.Equal(t, config.GetValueInt("key1.key2.intdata"), 12)
	assert.Equal(t, config.GetValueString("key1.key2.strdata"), "data")
	assert.Equal(t, config.GetValueBool("key1.key2.booldata"), true)
	assert.Equal(t, config.GetValueInt64("key1.key2.int64data"), isc.ToInt64(12))
	assert.Equal(t, config.GetValueFloat32("key1.key2.floatdata"), isc.ToFloat32(12.3))

	expectData := map[string]any{
		"field1": 1,
		"field2": "value2",
	}
	actData := map[string]any{}
	err := config.GetValueObject("key1.key2.objdata", &actData)
	if err != nil {
		return
	}
	assert.Equal(t, actData, expectData)
	assert.Equal(t, config.GetValueArrayInt("key1.key2.arraydata"), []int{1, 2, 3})
}

// 测试小驼峰
func TestSmall(t *testing.T) {
	config.LoadConfig()

	entity := SmallEntity{}
	config.GetValueObject("key1.ok", &entity)
	fmt.Println(entity)
}

type SmallEntity struct {
	HaoDeOk int
	NameAge int
}
