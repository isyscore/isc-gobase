package isc

type ISCMapToMap[K comparable, V comparable, R comparable] struct {
	m map[K]V
}

func MapToMapFrom[K comparable, V comparable, R comparable](m ISCMap[K, V]) ISCMapToMap[K, V, R] {
	return ISCMapToMap[K, V, R]{
		m: m.m,
	}
}

func (m ISCMapToMap[K, V, R]) FlatMap(f func(K, V) []R) ISCList[R] {
	r := MapFlatMap(m.m, f)
	return NewList(r)
}

func (m ISCMapToMap[K, V, R]) FlatMapTo(dest *[]R, f func(K, V) []R) ISCList[R] {
	r := MapFlatMapTo(m.m, dest, f)
	return NewList(r)
}

func (m ISCMapToMap[K, V, R]) Map(f func(K, V) R) ISCList[R] {
	r := MapMap(m.m, f)
	return NewList(r)
}

func (m ISCMapToMap[K, V, R]) MapTo(dest *[]R, f func(K, V) R) ISCList[R] {
	r := MapMapTo(m.m, dest, f)
	return NewList(r)
}
