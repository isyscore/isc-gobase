package isc

type ISCListToMap[T comparable, R comparable] struct {
	ISCList[T]
}

func ListToMapFrom[T comparable, R comparable](list ISCList[T]) ISCListToMap[T, R] {
	return ISCListToMap[T, R]{
		list,
	}
}

func (l ISCListToMap[T, R]) FlatMap(f func(T) []R) ISCList[R] {
	return ListFlatMap(l.ISCList, f)
}

func (l ISCListToMap[T, R]) FlatMapIndexed(f func(int, T) []R) ISCList[R] {
	return ListFlatMapIndexed(l.ISCList, f)
}

func (l ISCListToMap[T, R]) FlatMapTo(dest *[]R, f func(T) []R) ISCList[R] {
	return ListFlatMapTo(l.ISCList, dest, f)
}

func (l ISCListToMap[T, R]) FlatMapIndexedTo(dest *[]R, f func(int, T) []R) ISCList[R] {
	return ListFlatMapIndexedTo(l.ISCList, dest, f)
}

func (l ISCListToMap[T, R]) Map(f func(T) R) ISCList[R] {
	return ListMap(l.ISCList, f)
}

func (l ISCListToMap[T, R]) MapIndexed(f func(int, T) R) ISCList[R] {
	return ListMapIndexed(l.ISCList, f)
}

func (l ISCListToMap[T, R]) MapTo(dest *[]R, f func(T) R) ISCList[R] {
	return ListMapTo(l.ISCList, dest, f)
}

func (l ISCListToMap[T, R]) MapIndexedTo(dest *[]R, f func(int, T) R) ISCList[R] {
	return ListMapIndexedTo(l.ISCList, dest, f)
}

func (l ISCListToMap[T, R]) Reduce(init func(T) R, f func(R, T) R) R {
	return Reduce(l.ISCList, init, f)
}

func (l ISCListToMap[T, R]) ReduceIndexed(init func(int, T) R, f func(int, R, T) R) R {
	return ReduceIndexed(l.ISCList, init, f)
}

type ISCListToGroup[T comparable, K comparable, V comparable] struct {
	ISCList[T]
}

func ListToGroupFrom[T comparable, K comparable, V comparable](list ISCList[T]) ISCListToGroup[T, K, V] {
	return ISCListToGroup[T, K, V]{
		list,
	}
}

func (l ISCListToGroup[T, K, V]) GroupBy(f func(T) K) map[K][]T {
	return GroupBy(l.ISCList, f)
}

func (l ISCListToGroup[T, K, V]) GroupByTransform(f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransform(l.ISCList, f, trans)
}

func (l ISCListToGroup[T, K, V]) GroupByTo(dest *map[K][]T, f func(T) K) map[K][]T {
	return GroupByTo(l.ISCList, dest, f)
}

func (l ISCListToGroup[T, K, V]) GroupByTransformTo(dest *map[K][]V, f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransformTo(l.ISCList, dest, f, trans)
}

type ISCListToPair[K comparable, V comparable] struct {
	ISCList[Pair[K, V]]
}

func ListToPairFrom[K comparable, V comparable](list ISCList[Pair[K, V]]) ISCListToPair[K, V] {
	return ISCListToPair[K, V]{
		list,
	}
}

func ListToPairWithPairs[K comparable, V comparable](list ...Pair[K, V]) ISCListToPair[K, V] {
	return ISCListToPair[K, V]{
		list,
	}
}

func (l ISCListToPair[K, V]) ToMap() ISCMap[K, V] {
	m := make(map[K]V)
	for _, item := range l.ISCList {
		m[item.First] = item.Second
	}
	return NewMapWithMap(m)
}
