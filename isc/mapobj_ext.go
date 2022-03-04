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
	return MapFlatMap(m.ISCMap, f)
}

func (m ISCMapToMap[K, V, R]) FlatMapTo(dest *[]R, f func(K, V) []R) ISCList[R] {
	return MapFlatMapTo(m.ISCMap, dest, f)
}

func (m ISCMapToMap[K, V, R]) Map(f func(K, V) R) ISCList[R] {
	return MapMap(m.ISCMap, f)
}

func (m ISCMapToMap[K, V, R]) MapTo(dest *[]R, f func(K, V) R) ISCList[R] {
	return MapMapTo(m.ISCMap, dest, f)
}
