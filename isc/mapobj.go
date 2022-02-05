package isc

type ISCMap[K comparable, V comparable] struct {
	m map[K]V
}

func NewMap[K comparable, V comparable](ma map[K]V) ISCMap[K, V] {
	return ISCMap[K, V]{
		m: ma,
	}
}

func NewMapWithPairs[K comparable, V comparable](pairs ...Pair[K, V]) ISCMap[K, V] {
	m := make(map[K]V)
	for _, item := range pairs {
		m[item.First] = item.Second
	}
	return ISCMap[K, V]{
		m: m,
	}
}

func (m ISCMap[K, V]) ToMap() map[K]V {
	return m.m
}

func (m ISCMap[K, V]) Filter(f func(K, V) bool) ISCMap[K, V] {
	r := MapFilter(m.m, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) FilterNot(f func(K, V) bool) ISCMap[K, V] {
	r := MapFilterNot(m.m, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) FilterKeys(f func(K) bool) ISCMap[K, V] {
	r := MapFilterKeys(m.m, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) FilterValues(f func(V) bool) ISCMap[K, V] {
	r := MapFilterValues(m.m, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) FilterTo(dest *map[K]V, f func(K, V) bool) ISCMap[K, V] {
	r := MapFilterTo(m.m, dest, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) FilterNotTo(dest *map[K]V, f func(K, V) bool) ISCMap[K, V] {
	r := MapFilterNotTo(m.m, dest, f)
	return NewMap(r)
}

func (m ISCMap[K, V]) Contains(k K, v V) bool {
	return MapContains(m.m, k, v)
}

func (m ISCMap[K, V]) ContainsKey(k K) bool {
	return MapContainsKey(m.m, k)
}

func (m ISCMap[K, V]) ContainsValue(v V) bool {
	return MapContainsValue(m.m, v)
}

func (m ISCMap[K, V]) JoinToString(f func(K, V) string) string {
	return MapJoinToString(m.m, f)
}

func (m ISCMap[K, V]) JoinToStringFull(sep string, prefix string, postfix string, f func(K, V) string) string {
	return MapJoinToStringFull(m.m, sep, prefix, postfix, f)
}

func (m ISCMap[K, V]) All(f func(K, V) bool) bool {
	return MapAll(m.m, f)
}

func (m ISCMap[K, V]) Any(f func(K, V) bool) bool {
	return MapAny(m.m, f)
}

func (m ISCMap[K, V]) None(f func(K, V) bool) bool {
	return MapNone(m.m, f)
}

func (m ISCMap[K, V]) Count(f func(K, V) bool) int {
	return MapCount(m.m, f)
}

func (m ISCMap[K, V]) AllKey(f func(K) bool) bool {
	return MapAllKey(m.m, f)
}

func (m ISCMap[K, V]) AnyKey(f func(K) bool) bool {
	return MapAnyKey(m.m, f)
}

func (m ISCMap[K, V]) NoneKey(f func(K) bool) bool {
	return MapNoneKey(m.m, f)
}

func (m ISCMap[K, V]) CountKey(f func(K) bool) int {
	return MapCountKey(m.m, f)
}

func (m ISCMap[K, V]) AllValue(f func(V) bool) bool {
	return MapAllValue(m.m, f)
}

func (m ISCMap[K, V]) AnyValue(f func(V) bool) bool {
	return MapAnyValue(m.m, f)
}

func (m ISCMap[K, V]) NoneValue(f func(V) bool) bool {
	return MapNoneValue(m.m, f)
}

func (m ISCMap[K, V]) CountValue(f func(V) bool) int {
	return MapCountValue(m.m, f)
}

func (m ISCMap[K, V]) ToList() []Pair[K, V] {
	var n []Pair[K, V]
	for k, v := range m.m {
		n = append(n, NewPair(k, v))
	}
	return n
}
