package isc

import (
	"fmt"
)

type ISCSet[T any] ISCList[T]

// NewSet 初始化并指定存储对象的类型
func NewSet[T any]() ISCSet[T] {
	return ISCSet[T]{}
}

func NewSetWithList[T any](list []T) ISCSet[T] {
	s := ISCSet[T]{}
	s.AddAll(list...)
	return s
}

func NewSetWithItems[T any](items ...T) ISCSet[T] {
	s := ISCSet[T]{}
	s.AddAll(items...)
	return s
}

// Size 返回数据数量
func (s ISCSet[T]) Size() int {
	return len(s)
}

// Add 添加元素
func (s *ISCSet[T]) Add(item T) error {
	if !ISCList[T](*s).Contains(item) {
		*s = append(*s, item)
		return nil
	} else {
		return fmt.Errorf("%v already exists in set", item)
	}
}

// AddAll 添加多个元素
func (s *ISCSet[T]) AddAll(items ...T) {
	for _, item := range items {
		_ = s.Add(item)
	}
}

// Delete 删除指定Key元素
func (s *ISCSet[T]) Delete(item T) error {
	if idx := ISCList[T](*s).IndexOf(item); idx != -1 {
		*s = append((*s)[:idx], (*s)[idx+1:]...)
		return nil
	} else {
		return fmt.Errorf("%v not exists in set", item)
	}
}

// Contains 判断key是否存在
func (s ISCSet[T]) Contains(item T) bool {
	return ISCList[T](s).Contains(item)
}

// Clear 重置
func (s *ISCSet[T]) Clear() {
	*s = ISCSet[T]{}
}

func (s ISCSet[T]) ToList() ISCList[T] {
	return ISCList[T](s)
}
