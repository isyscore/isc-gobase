package isc

func ListToMap[K comparable, V any](list []Pair[K, V]) map[K]V {
	m := make(map[K]V)
	for _, item := range list {
		m[item.First] = item.Second
	}
	return m
}

func MapToList[K comparable, V any](m map[K]V) []Pair[K, V] {
	var n []Pair[K, V]
	for k, v := range m {
		n = append(n, NewPair(k, v))
	}
	return n
}
