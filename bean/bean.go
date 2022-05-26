package bean

import (
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"reflect"
	"strings"
)

var BeanMap map[string]any

func init() {
	BeanMap = map[string]any{}
}

// AddBean 添加bean
// 强烈建议：bean使用指针
func AddBean(beanName string, beanPtr any) {
	beanType := reflect.TypeOf(beanPtr)
	if beanType.Kind() != reflect.Ptr {
		logger.Warn("bean 只可为指针类型")
		return
	}
	BeanMap[beanName] = beanPtr
}

func GetBean(beanName string) any {
	if beanV, exit := BeanMap[beanName]; exit {
		return beanV
	}
	return nil
}

func Clean() {
	BeanMap = map[string]any{}
}

func GetBeanNames(beanName string) []string {
	if beanName == "" {
		j := 0
		keys := make([]string, len(BeanMap))
		for k := range BeanMap {
			keys[j] = k
			j++
		}
		return keys
	} else {
		j := 0
		keys := []string{}
		for k := range BeanMap {
			if strings.Contains(k, beanName) {
				keys = append(keys, k)
				j++
			}
		}
		return keys
	}
}

func ExistBean(beanName string) bool {
	_, exist := BeanMap[beanName]
	return exist
}

func CallFun(beanName, methodName string, parameterValueMap map[string]any) []any {
	if beanValue, exist := BeanMap[beanName]; exist {
		fType := reflect.TypeOf(beanValue)

		result := []any{}
		for index, num := 0, fType.NumMethod(); index < num; index++ {
			method := fType.Method(index)

			// 私有字段不处理
			if !isc.IsPublic(method.Name) {
				continue
			}

			if method.Name != methodName {
				continue
			}

			parameterNum := method.Type.NumIn()
			var in []reflect.Value
			in = append(in, reflect.ValueOf(beanValue))
			for i := 1; i < parameterNum; i++ {
				if v, exit := parameterValueMap["p"+isc.ToString(i)]; exit {
					in = append(in, reflect.ValueOf(v))
				}
			}

			if len(in) != parameterNum {
				return nil
			}
			vs := method.Func.Call(in)
			for _, v := range vs {
				result = append(result, v)
			}
		}
		return result
	}
	return nil
}

func GetField(beanName, fieldName string) any {
	if beanValue, exist := BeanMap[beanName]; exist {
		fValue := reflect.ValueOf(beanValue)
		fType := reflect.TypeOf(beanValue)

		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
			fValue = fValue.Elem()
		}

		for index, num := 0, fType.NumField(); index < num; index++ {
			field := fType.Field(index)

			// 私有字段不处理
			if !isc.IsPublic(field.Name) {
				continue
			}

			if field.Name == fieldName {
				return fValue.Field(index).Interface()
			}
		}
	}
	return nil
}

// SetField 修改属性的话，请将对象设置为指针，否则不生效
func SetField(beanName string, fieldName string, fieldValue any) {
	if beanValue, exist := BeanMap[beanName]; exist {
		// 私有字段不处理
		if !isc.IsPublic(fieldName) {
			return
		}

		fValue := reflect.ValueOf(beanValue)
		fType := reflect.TypeOf(beanValue)

		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
			fValue = fValue.Elem()
		} else {
			return
		}

		if _, exist := fType.FieldByName(fieldName); exist {
			fValue.FieldByName(fieldName).Set(reflect.ValueOf(fieldValue))
		}
	}
}
