package test

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/magiconair/properties/assert"
	"os"
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

	entity1 := SmallEntity{}
	_ = config.GetValueObject("key1.ok1", &entity1)
	assert.Equal(t, entity1.NameAge, 32)
	assert.Equal(t, entity1.HaoDeOk, 12)

	entity2 := SmallEntity{}
	_ = config.GetValueObject("key1.ok2", &entity2)
	assert.Equal(t, entity2.NameAge, 32)
	assert.Equal(t, entity2.HaoDeOk, 12)

	entity3 := SmallEntity{}
	_ = config.GetValueObject("key1.ok3", &entity3)
	assert.Equal(t, entity3.NameAge, 32)
	assert.Equal(t, entity3.HaoDeOk, 12)

	entity4 := SmallEntity{}
	_ = config.GetValueObject("key1.ok4", &entity4)
	assert.Equal(t, entity4.NameAge, 32)
	assert.Equal(t, entity4.HaoDeOk, 12)
}

type SmallEntity struct {
	HaoDeOk int
	NameAge int
}

// 测试兼容标签：json、yaml
func TestJsonOrYamlTag(t *testing.T) {
	config.LoadConfig()

	entity1 := SmallEntityJsonTag{}
	_ = config.GetValueObject("key1.json", &entity1)
	assert.Equal(t, entity1.NameAge, 32)
	assert.Equal(t, entity1.HaoDeOk, 12)

	entity1_1 := SmallEntityJsonTag2{}
	_ = config.GetValueObject("key1.json", &entity1_1)
	assert.Equal(t, entity1_1.NameAge, 32)
	assert.Equal(t, entity1_1.HaoDeOk, 12)

	entity2 := SmallEntityYamlTag{}
	_ = config.GetValueObject("key1.yaml", &entity2)
	assert.Equal(t, entity2.NameAge, 32)
	assert.Equal(t, entity2.HaoDeOk, 12)

	entity2_1 := SmallEntityYamlTag2{}
	_ = config.GetValueObject("key1.yaml", &entity2_1)
	assert.Equal(t, entity2_1.NameAge, 32)
	assert.Equal(t, entity2_1.HaoDeOk, 12)
}

type SmallEntityJsonTag struct {
	HaoDeOk int `json:"test_haode"`
	NameAge int `json:"test_namehaha"`
}

type SmallEntityJsonTag2 struct {
	HaoDeOk int `json:"test_haode,omitempty"`
	NameAge int `json:"test_namehaha"`
}

type SmallEntityYamlTag struct {
	HaoDeOk int `yaml:"test_haode"`
	NameAge int `yaml:"test_namehaha"`
}

type SmallEntityYamlTag2 struct {
	HaoDeOk int `yaml:"test_haode,omitempty"`
	NameAge int `yaml:"test_namehaha,flow"`
}

// 测试读取某个文件
func TestRead(t *testing.T) {
	config.LoadFile("./application-local.yaml")

	en := EntityTest{}
	err := config.GetValueObject("entity", &en)
	if err != nil {
		logger.Warn("转换告警")
		return
	}

	assert.Equal(t, en.Name, "name-change")
}

type EntityTest struct {
	Name string
}

// 测试append
func TestAppend(t *testing.T) {
	config.LoadFile("./application-append-original.yaml")
	config.AppendFile("./application-append.yaml")

	assert.Equal(t, config.GetValueString("a.b.c"), "c-value-change")
	assert.Equal(t, config.GetValueString("a.b.d"), "d-value")
	assert.Equal(t, config.GetValueString("a.b.e.f"), "f-value")
}

// 测试cm文件位置定制化
func TestConfigInit(t *testing.T) {
	_ = os.Setenv("base.config.additional-location", "./application-append.yaml")
	config.LoadConfigFromRelativePath("./application-append-original.yaml")

	assert.Equal(t, config.GetValueString("a.b.c"), "c-value-change")
	assert.Equal(t, config.GetValueString("a.b.d"), "d-value")
	assert.Equal(t, config.GetValueString("a.b.e.f"), "f-value")
}

// 测试：yaml占位符的功能
func TestPlaceHolder1(t *testing.T) {
	config.LoadFile("./application-place1.yml")

	assert.Equal(t, config.GetValueString("place.name"), "test")
	assert.Equal(t, config.GetValueString("test.name"), "test")
	assert.Equal(t, config.GetValueString("test.name2"), "test2")
}

// 测试：yaml占位符的功能
func TestPlaceHolder2(t *testing.T) {
	config.LoadFile("./application-place1.yaml")

	assert.Equal(t, config.GetValueString("place.name"), "test")
	assert.Equal(t, config.GetValueString("test.name"), "test")
	assert.Equal(t, config.GetValueString("test.name2"), "test2")
}

// 测试：yaml占位符的功能
func TestPlaceHolder3(t *testing.T) {
	config.LoadFile("./application-place1.json")

	assert.Equal(t, config.GetValueString("place.name"), "test")
	assert.Equal(t, config.GetValueString("test.name"), "test")
	assert.Equal(t, config.GetValueString("test.name2"), "test2")
}

// 测试：yaml占位符的功能
func TestPlaceHolder4(t *testing.T) {
	config.LoadFile("./application-place1.properties")

	assert.Equal(t, config.GetValueString("place.name"), "test")
	assert.Equal(t, config.GetValueString("test.name"), "test")
	assert.Equal(t, config.GetValueString("test.name2"), "test2")
}
