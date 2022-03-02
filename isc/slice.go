package isc

//SubList 分片截取
func SubList[T any](list []T, fromIndex int, toIndex int) []T {
	m := map[int]int{fromIndex: toIndex, toIndex: fromIndex}
	start := fromIndex
	if fromIndex < start {
		start = toIndex
	}
	end := m[start]

	if start < 0 {
		start = 0
	}
	if start == end {
		return []T{}
	}
	if end > len(list) {
		return list[start:]
	}
	return list[start:end]
}

//Slice 分片截取,参数详情见 IntRange ,返回新分片
func Slice[T any](list []T, r IntRange) []T {
	return SubList(list, r.Start, r.End)
}

// SliceContains Returns true if element is found in the collection.
//predicate keySelector
//if you want to check item in list, please use ListContains
func SliceContains[T any, K comparable](list []T, predicate func(T) K, key K) bool {
	m := SliceTo(list, predicate)
	_, ok := m[key]
	return ok
}

func IsInSlice[T comparable](list []T, val T) bool {
	return IndexOf(list, val) >= 0
}

func SliceToMap[T comparable](list []T) map[T]T {
	m := make(map[T]T)
	for _, e := range list {
		m[e] = e
	}
	return m
}

func SliceTo[T any, K comparable](list []T, valueTransform func(T) K) map[K]T {
	m := make(map[K]T)
	for _, e := range list {
		m[valueTransform(e)] = e
	}
	return m
}

func SliceDistinct[T any](list []T) []T {
	result := NewList[T]()
	for _, k := range list {
		if !result.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

//SliceDistinctTo Returns a list containing only distinct elements from the given collection.
//Among equal elements of the given collection, only the last one will be present in the resulting list.
//The elements in the resulting list are not in the same order as they were in the source collection.
func SliceDistinctTo[T any, V comparable](list []T, valueTransform func(T) V) []T {
	m := SliceTo(list, valueTransform)
	var result []T
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
