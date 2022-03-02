package isc

import (
	"reflect"
)

//ListFilter filter specificated item in a list
func ListFilter[T any](list []T, f func(T) bool) []T {
	var dest []T
	return ListFilterTo(list, &dest, f)
}

//ListFilterNot Returns a list containing all elements not matching the given predicate.
func ListFilterNot[T any](list []T, predicate func(T) bool) []T {
	var n []T
	return ListFilterNotTo(list, &n, predicate)
}

//ListFilterIndexed Returns a list containing only elements matching the given predicate.
//Params: predicate - function that takes the index of an element and the element itself and returns the result of predicate evaluation on the element.
func ListFilterIndexed[T any](list []T, predicate func(int, T) bool) []T {
	var n []T
	return ListFilterIndexedTo(list, &n, predicate)
}

//ListFilterNotIndexed Appends all elements matching the given predicate to the given destination.
//Params: predicate - function that takes the index of an element and the element itself and returns the result of predicate evaluation on the element.
func ListFilterNotIndexed[T any](list []T, f func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if !f(i, e) {
			n = append(n, e)
		}
	}
	return n
}

//ListFilterNotNull Returns a list containing all elements that are not null.
func ListFilterNotNull[T any](list []*T) []*T {
	var n []*T
	for _, e := range list {
		if e != nil {
			n = append(n, e)
		}
	}
	return n
}

//ListFilterTo Appends all elements matching the given predicate to the given dest.
func ListFilterTo[T any](list []T, dest *[]T, f func(T) bool) []T {
	var n []T
	for _, e := range list {
		if f(e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

//ListFilterNotTo Appends all elements not matching the given predicate to the given destination.
func ListFilterNotTo[T any](list []T, dest *[]T, predicate func(T) bool) []T {
	var n []T
	for _, e := range list {
		if !predicate(e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

//ListFilterIndexedTo Appends all elements matching the given predicate to the given destination.
//Params: predicate - function that takes the index of an element and the element itself and returns the result of predicate evaluation on the element.
func ListFilterIndexedTo[T any](list []T, dest *[]T, predicate func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if predicate(i, e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

//ListFilterNotIndexedTo Appends all elements not matching the given predicate to the given destination.
//Params: predicate - function that takes the index of an element and the element itself and returns the result of predicate evaluation on the element.
func ListFilterNotIndexedTo[T any](list []T, dest *[]T, predicate func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if !predicate(i, e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

func ListContains[T any](list []T, item T) bool {
	ret := false
	for _, e := range list {
		if reflect.DeepEqual(e, item) {
			ret = true
			break
		}
	}
	return ret
}

func ListDistinct[T any](list []T) []T {
	return SliceDistinct(list)
}

/// functions for map

func MapFilter[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterNot[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if !f(k, v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if f(k) {
			n[k] = v
		}
	}
	return n
}

func MapFilterValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if f(v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterTo[K comparable, V any](m map[K]V, dest *map[K]V, f func(K, V) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			(*dest)[k] = v
			n[k] = v
		}
	}
	return n
}

func MapFilterNotTo[K comparable, V any](m map[K]V, dest *map[K]V, f func(K, V) bool) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		if !f(k, v) {
			(*dest)[k] = v
			n[k] = v
		}
	}
	return n
}

func MapContains[K comparable, V any](m map[K]V, k K, v V) bool {
	ret := false
	for t, u := range m {
		if t == k && reflect.DeepEqual(u, v) {
			ret = true
			break
		}
	}
	return ret
}

func MapContainsKey[K comparable, V any](m map[K]V, k K) bool {
	ret := false
	for t := range m {
		if t == k {
			ret = true
			break
		}
	}
	return ret
}

func MapContainsValue[K comparable, V any](m map[K]V, v V) bool {
	ret := false
	for _, u := range m {
		if reflect.DeepEqual(u, v) {
			ret = true
			break
		}
	}
	return ret
}
