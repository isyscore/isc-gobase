package isc

type ISCMap[K comparable, V any] map[K]V

func NewMap[K comparable, V any]() ISCMap[K, V] {
	return make(ISCMap[K, V])
}

func NewMapWithMap[K comparable, V any](ma map[K]V) ISCMap[K, V] {
	return ma
}

func NewMapWithPairs[K comparable, V any](pairs ...Pair[K, V]) ISCMap[K, V] {
	m := make(map[K]V)
	for _, item := range pairs {
		m[item.First] = item.Second
	}
	return m
}

func (m ISCMap[K, V]) Size() int {
	return len(m)
}

func (m ISCMap[K, V]) Put(k K, v V) {
	m[k] = v
}

func (m ISCMap[K, V]) PutPair(item Pair[K, V]) {
	m[item.First] = item.Second
}

func (m ISCMap[K, V]) PutAllPairs(item ...Pair[K, V]) {
	for _, e := range item {
		m.PutPair(e)
	}
}

func (m ISCMap[K, V]) Get(k K) V {
	return m[k]
}

func (m ISCMap[K, V]) GetOrDef(k K, def V) V {
	if v, ok := m[k]; ok {
		return v
	} else {
		return def
	}
}

func (m ISCMap[K, V]) Delete(k K) {
	delete(m, k)
}

func (m *ISCMap[K, V]) Clear() {
	*m = make(ISCMap[K, V])
}

func (m ISCMap[K, V]) ForEach(f func(K, V)) {
	for k, v := range m {
		f(k, v)
	}
}

func (m ISCMap[K, V]) Filter(f func(K, V) bool) ISCMap[K, V] {
	return MapFilter(m, f)
}

func (m ISCMap[K, V]) FilterNot(f func(K, V) bool) ISCMap[K, V] {
	return MapFilterNot(m, f)
}

func (m ISCMap[K, V]) FilterKeys(f func(K) bool) ISCMap[K, V] {
	return MapFilterKeys(m, f)
}

func (m ISCMap[K, V]) FilterValues(f func(V) bool) ISCMap[K, V] {
	return MapFilterValues(m, f)
}

func (m ISCMap[K, V]) FilterTo(dest *map[K]V, f func(K, V) bool) ISCMap[K, V] {
	return MapFilterTo(m, dest, f)
}

func (m ISCMap[K, V]) FilterNotTo(dest *map[K]V, f func(K, V) bool) ISCMap[K, V] {
	return MapFilterNotTo(m, dest, f)
}

func (m ISCMap[K, V]) Contains(k K, v V) bool {
	return MapContains(m, k, v)
}

func (m ISCMap[K, V]) ContainsKey(k K) bool {
	return MapContainsKey(m, k)
}

func (m ISCMap[K, V]) ContainsValue(v V) bool {
	return MapContainsValue(m, v)
}

func (m ISCMap[K, V]) JoinToString(f func(K, V) string) string {
	return MapJoinToString(m, f)
}

func (m ISCMap[K, V]) JoinToStringFull(sep string, prefix string, postfix string, f func(K, V) string) string {
	return MapJoinToStringFull(m, sep, prefix, postfix, f)
}

func (m ISCMap[K, V]) All(f func(K, V) bool) bool {
	return MapAll(m, f)
}

func (m ISCMap[K, V]) Any(f func(K, V) bool) bool {
	return MapAny(m, f)
}

func (m ISCMap[K, V]) None(f func(K, V) bool) bool {
	return MapNone(m, f)
}

func (m ISCMap[K, V]) Count(f func(K, V) bool) int {
	return MapCount(m, f)
}

func (m ISCMap[K, V]) AllKey(f func(K) bool) bool {
	return MapAllKey(m, f)
}

func (m ISCMap[K, V]) AnyKey(f func(K) bool) bool {
	return MapAnyKey(m, f)
}

func (m ISCMap[K, V]) NoneKey(f func(K) bool) bool {
	return MapNoneKey(m, f)
}

func (m ISCMap[K, V]) CountKey(f func(K) bool) int {
	return MapCountKey(m, f)
}

func (m ISCMap[K, V]) AllValue(f func(V) bool) bool {
	return MapAllValue(m, f)
}

func (m ISCMap[K, V]) AnyValue(f func(V) bool) bool {
	return MapAnyValue(m, f)
}

func (m ISCMap[K, V]) NoneValue(f func(V) bool) bool {
	return MapNoneValue(m, f)
}

func (m ISCMap[K, V]) CountValue(f func(V) bool) int {
	return MapCountValue(m, f)
}

func (m ISCMap[K, V]) ToList() []Pair[K, V] {
	var n []Pair[K, V]
	for k, v := range m {
		n = append(n, NewPair(k, v))
	}
	return n
}

func (m ISCMap[K, V]) Plus(n map[K]V) ISCMap[K, V] {
	return MapPlus(m, n)
}

func (m ISCMap[K, V]) Minus(n map[K]V) ISCMap[K, V] {
	return MapMinus(m, n)
}

func (m ISCMap[K, V]) Equals(n map[K]V) bool {
	return MapEquals(m, n)
}

func (m ISCMap[K, V]) Keys() ISCList[K] {
	i := 0
	keys := make([]K, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
