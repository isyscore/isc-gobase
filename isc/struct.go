package isc

import (
	"reflect"
	"unsafe"
)

// 获取对象的属性：一般用于访问私有属性
func GetPrivateFieldValue(objPtrValue reflect.Value, fieldName string) interface{} {
	if objPtrValue.Kind() != reflect.Ptr {
		return nil
	}
	fieldValue := objPtrValue.Elem().FieldByName(fieldName)
	return reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem().Interface()
}

// 给对象的属性设置值：一般用于设置私有属性
func SetFieldPrivateValue(objPtrValue reflect.Value, fieldName string, fieldNewValue reflect.Value) {
	if objPtrValue.Kind() != reflect.Ptr {
		return
	}
	fieldValue := objPtrValue.Elem().FieldByName(fieldName)
	fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
	fieldValue.Set(fieldNewValue.Elem())
}
