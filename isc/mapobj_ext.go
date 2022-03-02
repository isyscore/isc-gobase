package isc

type ISCMapToMap[K comparable, V any, R any] struct {
	ISCMap[K, V]
}

func MapToMapFrom[K comparable, V any, R any](m ISCMap[K, V]) ISCMapToMap[K, V, R] {
	return ISCMapToMap[K, V, R]{
		m,
	}
}

func (m ISCMapToMap[K, V, R]) FlatMap(f func(K, V) []R) ISCList[R] {
	r := MapFlatMap(m.ISCMap, f)
	return NewListWithList(r)
}

func (m ISCMapToMap[K, V, R]) FlatMapTo(dest *[]R, f func(K, V) []R) ISCList[R] {
	r := MapFlatMapTo(m.ISCMap, dest, f)
	return NewListWithList(r)
}

func (m ISCMapToMap[K, V, R]) Map(f func(K, V) R) ISCList[R] {
	r := MapMap(m.ISCMap, f)
	return NewListWithList(r)
}

func (m ISCMapToMap[K, V, R]) MapTo(dest *[]R, f func(K, V) R) ISCList[R] {
	r := MapMapTo(m.ISCMap, dest, f)
	return NewListWithList(r)
}
