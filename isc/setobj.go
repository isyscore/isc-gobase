package isc

import (
	"fmt"
)

type ISCSet[T comparable] map[T]struct{}

// NewSet 初始化并指定存储对象的类型
func NewSet[T comparable]() ISCSet[T] {
	return ISCSet[T]{}
}

func NewSetWithList[T comparable](list []T) ISCSet[T] {
	s := ISCSet[T]{}
	s.AddAll(list...)
	return s
}

func NewSetWithItems[T comparable](items ...T) ISCSet[T] {
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
	if !s.Contains(item) {
		(*s)[item] = struct{}{}
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
	if s.Contains(item) {
		delete(*s, item)
		return nil
	} else {
		return fmt.Errorf("%v not exists in set", item)
	}
}

// Contains 判断key是否存在
func (s ISCSet[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Clear 重置
func (s *ISCSet[T]) Clear() {
	*s = ISCSet[T]{}
}

func (s ISCSet[T]) ToList() ISCList[T] {
	res := NewList[T]()
	for k := range s {
		res.Add(k)
	}
	return res
}
