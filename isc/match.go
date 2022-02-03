package isc

func ListAll[T any](list []T, f func(T) bool) bool {
	for _, e := range list {
		if !f(e) {
			return false
		}
	}
	return true
}

func ListAny[T any](list []T, f func(T) bool) bool {
	for _, e := range list {
		if f(e) {
			return true
		}
	}
	return false
}

func ListNone[T any](list []T, f func(T) bool) bool {
	for _, e := range list {
		if f(e) {
			return false
		}
	}
	return true
}

func ListCount[T any](list []T, f func(T) bool) int {
	num := 0
	for _, e := range list {
		if f(e) {
			num++
		}
	}
	return num
}

func MapAll[K comparable, V any](m map[K]V, f func(K, V) bool) bool {
	for k, v := range m {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func MapAny[K comparable, V any](m map[K]V, f func(K, V) bool) bool {
	for k, v := range m {
		if f(k, v) {
			return true
		}
	}
	return false
}

func MapNone[K comparable, V any](m map[K]V, f func(K, V) bool) bool {
	for k, v := range m {
		if f(k, v) {
			return false
		}
	}
	return true
}

func MapCount[K comparable, V any](m map[K]V, f func(K, V) bool) int {
	num := 0
	for k, v := range m {
		if f(k, v) {
			num++
		}
	}
	return num
}

func MapAllKey[K comparable, V any](m map[K]V, f func(K) bool) bool {
	for k := range m {
		if !f(k) {
			return false
		}
	}
	return true
}

func MapAnyKey[K comparable, V any](m map[K]V, f func(K) bool) bool {
	for k := range m {
		if f(k) {
			return true
		}
	}
	return false
}

func MapNoneKey[K comparable, V any](m map[K]V, f func(K) bool) bool {
	for k := range m {
		if f(k) {
			return false
		}
	}
	return true
}

func MapCountKey[K comparable, V any](m map[K]V, f func(K) bool) int {
	num := 0
	for k := range m {
		if f(k) {
			num++
		}
	}
	return num
}

func MapAllValue[K comparable, V any](m map[K]V, f func(V) bool) bool {
	for _, v := range m {
		if !f(v) {
			return false
		}
	}
	return true
}

func MapAnyValue[K comparable, V any](m map[K]V, f func(V) bool) bool {
	for _, v := range m {
		if f(v) {
			return true
		}
	}
	return false
}

func MapNoneValue[K comparable, V any](m map[K]V, f func(V) bool) bool {
	for _, v := range m {
		if f(v) {
			return false
		}
	}
	return true
}

func MapCountValue[K comparable, V any](m map[K]V, f func(V) bool) int {
	num := 0
	for _, v := range m {
		if f(v) {
			num++
		}
	}
	return num
}
