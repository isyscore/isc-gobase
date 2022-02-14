package isc

// filter specificated item in a list

func ListFilter[T any](list []T, f func(T) bool) []T {
	var n []T
	for _, e := range list {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func ListFilterNot[T any](list []T, f func(T) bool) []T {
	var n []T
	for _, e := range list {
		if !f(e) {
			n = append(n, e)
		}
	}
	return n
}

func ListFilterIndexed[T any](list []T, f func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if f(i, e) {
			n = append(n, e)
		}
	}
	return n
}

func ListFilterNotIndexed[T any](list []T, f func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if !f(i, e) {
			n = append(n, e)
		}
	}
	return n
}

func ListFilterNotNull[T any](list []*T) []*T {
	var n []*T
	for _, e := range list {
		if e != nil {
			n = append(n, e)
		}
	}
	return n
}

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

func ListFilterNotTo[T any](list []T, dest *[]T, f func(T) bool) []T {
	var n []T
	for _, e := range list {
		if !f(e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

func ListFilterIndexedTo[T any](list []T, dest *[]T, f func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if f(i, e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

func ListFilterNotIndexedTo[T any](list []T, dest *[]T, f func(int, T) bool) []T {
	var n []T
	for i, e := range list {
		if !f(i, e) {
			*dest = append(*dest, e)
			n = append(n, e)
		}
	}
	return n
}

func ListContains[T comparable](list []T, item T) bool {
	var ret = false
	for _, e := range list {
		if e == item {
			ret = true
			break
		}
	}
	return ret
}

func ListDistinct[T comparable](list []T) []T {
	return SliceDistinct(list)
}

/// functions for map

func MapFilter[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterNot[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if !f(k, v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if f(k) {
			n[k] = v
		}
	}
	return n
}

func MapFilterValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if f(v) {
			n[k] = v
		}
	}
	return n
}

func MapFilterTo[K comparable, V any](m map[K]V, dest *map[K]V, f func(K, V) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			(*dest)[k] = v
			n[k] = v
		}
	}
	return n
}

func MapFilterNotTo[K comparable, V any](m map[K]V, dest *map[K]V, f func(K, V) bool) map[K]V {
	var n map[K]V = make(map[K]V)
	for k, v := range m {
		if !f(k, v) {
			(*dest)[k] = v
			n[k] = v
		}
	}
	return n
}

func MapContains[K comparable, V comparable](m map[K]V, k K, v V) bool {
	var ret = false
	for t, u := range m {
		if t == k && u == v {
			ret = true
			break
		}
	}
	return ret
}

func MapContainsKey[K comparable, V any](m map[K]V, k K) bool {
	var ret = false
	for t := range m {
		if t == k {
			ret = true
			break
		}
	}
	return ret
}

func MapContainsValue[K comparable, V comparable](m map[K]V, v V) bool {
	var ret = false
	for _, u := range m {
		if u == v {
			ret = true
			break
		}
	}
	return ret
}
