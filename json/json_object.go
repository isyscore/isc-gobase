package json

import (
	"fmt"
	"github.com/isyscore/isc-gobase/isc"
	"reflect"
	"strings"
)

type Object struct {
	ValueMap     map[string]any
	ValueDeepMap map[string]any
}

func (jsonObject *Object) Load(jsonContent string) error {
	if jsonObject.ValueMap == nil {
		jsonObject.ValueMap = make(map[string]any)
	}

	if jsonObject.ValueDeepMap == nil {
		jsonObject.ValueDeepMap = make(map[string]any)
	}

	yamlStr, _ := isc.JsonToYaml(jsonContent)
	property, _ := isc.YamlToProperties(yamlStr)
	valueMap, _ := isc.PropertiesToMap(property)
	jsonObject.ValueMap = valueMap

	yamlMap, _ := isc.YamlToMap(yamlStr)
	jsonObject.ValueDeepMap = yamlMap
	return nil
}

func (jsonObject *Object) Put(key string, value any) {
	if nil == value {
		return
	}

	if oldValue, exist := jsonObject.ValueMap[key]; exist {
		if reflect.TypeOf(oldValue) != reflect.TypeOf(oldValue) {
			return
		}
	}
	jsonObject.ValueMap[key] = value
	jsonObject.doPut(key, value)
}

func (jsonObject *Object) doPut(key string, value any) {
	innerPutValue(jsonObject.ValueDeepMap, key, value)
}

func innerPutValue(valueMap map[string]any, key string, newValue any) {
	if nil == valueMap {
		valueMap = make(map[string]any)
	}
	if !strings.Contains(key, ".") {
		if oldValue, exist := valueMap[key]; exist {
			if reflect.TypeOf(oldValue) != reflect.TypeOf(newValue) {
				return
			}
		}
		valueMap[key] = newValue
		return
	}

	lastIndex := strings.Index(key, ".")
	startKey := key[:lastIndex]
	endKey := key[lastIndex+1:]
	if oldValue, exist := valueMap[startKey]; exist {
		if reflect.TypeOf(oldValue).Kind() != reflect.Map {
			return
		} else {
			oldValueMap := isc.ToMap(oldValue)
			innerPutValue(oldValueMap, endKey, newValue)
			valueMap[startKey] = oldValueMap
		}
	} else {
		oldValueMap := make(map[string]any)
		innerPutValue(oldValueMap, endKey, newValue)
		valueMap[startKey] = oldValueMap
	}
}

func (jsonObject *Object) Get(key string) any {
	return jsonObject.doGet(jsonObject.ValueDeepMap, key)
}

func (jsonObject *Object) initValue(key string) {
	jsonObject.doInitValue(jsonObject.ValueDeepMap, key)
}

func (jsonObject *Object) doInitValue(parentValue any, key string) {
	if key == "" {
		return
	}
	parentValueKind := reflect.ValueOf(parentValue).Kind()
	if parentValueKind == reflect.Map {
		keys := strings.SplitN(key, ".", 2)
		v1 := reflect.ValueOf(parentValue).MapIndex(reflect.ValueOf(keys[0]))
		emptyValue := reflect.Value{}
		if v1 == emptyValue {
			return
		}
		if len(keys) == 1 {
			jsonObject.doInitValue(v1.Interface(), "")
		} else {
			jsonObject.doInitValue(v1.Interface(), fmt.Sprintf("%v", keys[1]))
		}
	}
	return
}

func (jsonObject *Object) doGet(parentValue any, key string) any {
	if key == "" {
		return parentValue
	}
	parentValueKind := reflect.ValueOf(parentValue).Kind()
	if parentValueKind == reflect.Map {
		keys := strings.SplitN(key, ".", 2)

		currentKey := keys[0]
		nextKey := keys[1]
		if strings.Contains(keys[0], "[") && strings.Contains(keys[0], "]") {
			startIndex := strings.Index(keys[0], "[")
			endIndex := strings.Index(keys[0], "]")
			dataIndex := keys[0][startIndex+1 : endIndex]
			if isc.ToInt(dataIndex) >= 0 {
				currentKey = keys[0][:startIndex]
				nextKey = key[startIndex:]
			}
		}
		v1 := reflect.ValueOf(parentValue).MapIndex(reflect.ValueOf(currentKey))
		emptyValue := reflect.Value{}
		if v1 == emptyValue {
			return nil
		}
		if len(keys) == 1 {
			return v1.Interface()
		} else {
			return jsonObject.doGet(v1.Interface(), fmt.Sprintf("%v", nextKey))
		}
	} else if parentValueKind == reflect.Slice {
		// todo 柿子 还没写完，继续
		//if !strings.Contains(key, "[") && !strings.Contains(key, "]") {
		//	return nil
		//}
		//
		//startIndex := strings.Index(key, "[")
		//endIndex := strings.Index(key, "]")
		//dataIndex := isc.ToInt(key[startIndex+1 : endIndex])
		//
		//parentValueT := reflect.ValueOf(parentValue)
		//arrayLen := parentValueT.Len()
		//for arrayIndex := 0; arrayIndex < arrayLen; arrayIndex++ {
		//	if arrayIndex == dataIndex {
		//		fieldValue := parentValueT.Index(arrayIndex)
		//
		//		currentKey := keys[0]
		//		nextKey := keys[1]
		//		if strings.Contains(keys[0], "[") && strings.Contains(keys[0], "]") {
		//			startIndex := strings.Index(keys[0], "[")
		//			endIndex := strings.Index(keys[0], "]")
		//			dataIndex := keys[0][startIndex+1 : endIndex]
		//			if isc.ToInt(dataIndex) >= 0 {
		//				currentKey = keys[0][:startIndex]
		//				nextKey = key[startIndex:]
		//			}
		//		}
		//		v1 := reflect.ValueOf(parentValue).MapIndex(reflect.ValueOf(currentKey))
		//		emptyValue := reflect.Value{}
		//		if v1 == emptyValue {
		//			return nil
		//		}
		//		if len(keys) == 1 {
		//			return v1.Interface()
		//		} else {
		//			return jsonObject.doGet(v1.Interface(), fmt.Sprintf("%v", nextKey))
		//		}
		//	}
		//}
	}
	return nil
}

func (jsonObject *Object) GetString(key string) string {
	return isc.ToString(jsonObject.Get(key))
}

func (jsonObject *Object) GetInt(key string) int {
	return isc.ToInt(jsonObject.Get(key))
}

func (jsonObject *Object) GetInt8(key string) int8 {
	return isc.ToInt8(jsonObject.Get(key))
}

func (jsonObject *Object) GetInt16(key string) int16 {
	return isc.ToInt16(jsonObject.Get(key))
}

func (jsonObject *Object) GetInt32(key string) int32 {
	return isc.ToInt32(jsonObject.Get(key))
}

func (jsonObject *Object) GetInt64(key string) int64 {
	return isc.ToInt64(jsonObject.Get(key))
}

func (jsonObject *Object) GetUInt(key string) uint {
	return isc.ToUInt(jsonObject.Get(key))
}

func (jsonObject *Object) GetUInt8(key string) uint8 {
	return isc.ToUInt8(jsonObject.Get(key))
}

func (jsonObject *Object) GetUInt16(key string) uint16 {
	return isc.ToUInt16(jsonObject.Get(key))
}

func (jsonObject *Object) GetUInt32(key string) uint32 {
	return isc.ToUInt32(jsonObject.Get(key))
}

func (jsonObject *Object) GetUInt64(key string) uint64 {
	return isc.ToUInt64(jsonObject.Get(key))
}

func (jsonObject *Object) GetFloat32(key string) float32 {
	return isc.ToFloat32(jsonObject.Get(key))
}

func (jsonObject *Object) GetFloat64(key string) float64 {
	return isc.ToFloat64(jsonObject.Get(key))
}

func (jsonObject *Object) GetBool(key string) bool {
	return isc.ToBool(jsonObject.Get(key))
}

func (jsonObject *Object) GetObject(key string, targetPtrObj any) error {
	data := jsonObject.Get(key)
	err := isc.DataToObject(data, targetPtrObj)
	if err != nil {
		return err
	}
	return nil
}

func (jsonObject *Object) GetArray(key string) []any {
	var arrayResult = []any{}
	data := jsonObject.Get(key)
	err := isc.DataToObject(data, &arrayResult)
	if err != nil {
		return arrayResult
	}
	return arrayResult
}
