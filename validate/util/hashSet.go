package util

import (
	"errors"
	"reflect"
)

type HashSet struct {
	//数据载体
	data map[interface{}]interface{}
	//数据类型
	dataType reflect.Kind
	//数据数量
	count int
}

// NewHashSet 初始化并指定存储对象的类型
func NewHashSet(kind reflect.Kind) *HashSet {
	hashSet := new(HashSet)
	hashSet.data = make(map[interface{}]interface{})
	hashSet.dataType = kind
	return hashSet
}

// Size 返回数据数量
func (hashSet *HashSet) Size() int {
	return hashSet.count
}

// GetDataKind 返回数据类型
func (hashSet *HashSet) GetDataKind() interface{} {
	return hashSet.dataType
}

// Add 添加元素
func (hashSet *HashSet) Add(key interface{}) error {
	err := hashSet.checkData(key)
	if err != nil {
		return err
	}

	_, ok := hashSet.data[key]
	if ok {
		return errors.New("DataIsExist")
	}
	hashSet.count += 1
	hashSet.data[key] = key
	return nil
}

// AddMulti 添加多个元素
func (hashSet *HashSet) AddMulti(keyS ...interface{}) {
	for _, key := range keyS {
		err := hashSet.Add(key)
		if err != nil {
			continue
		}
	}
}

// Remove 删除指定Key元素
func (hashSet *HashSet) Remove(key interface{}) error {
	err := hashSet.checkData(key)
	if err != nil {
		return err
	}

	value, ok := hashSet.data[key]
	if ok {
		delete(hashSet.data, value)
		hashSet.count -= 1
		return nil
	}
	return errors.New("NotFoundKey")
}

// Contains 判断key是否存在
func (hashSet *HashSet) Contains(key interface{}) bool {
	err := hashSet.checkData(key)
	if err != nil {
		return false
	}
	_, ok := hashSet.data[key]
	if ok {
		return true
	} else {
		return false
	}
}

// Clear 重置
func (hashSet *HashSet) Clear() {
	hashSet.count = 0
	hashSet.data = make(map[interface{}]interface{})
}

// 判断添加元素是否为指定类型
func (hashSet *HashSet) checkData(data interface{}) error {
	if data == nil {
		return errors.New("dataIsNil")
	}

	dataKind := reflect.TypeOf(data).Kind()
	if hashSet.dataType != dataKind {
		return errors.New("UnsupportedTypes")
	}
	return nil
}
