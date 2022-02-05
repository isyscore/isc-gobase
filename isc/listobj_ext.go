package isc

type ISCListToMap[T comparable, R comparable] struct {
	l []T
}

func ListToMapFrom[T comparable, R comparable](list ISCList[T]) ISCListToMap[T, R] {
	return ISCListToMap[T, R]{
		l: list.l,
	}
}

func (l ISCListToMap[T, R]) FlatMap(f func(T) []R) ISCList[R] {
	r := ListFlatMap(l.l, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) FlatMapIndexed(f func(int, T) []R) ISCList[R] {
	r := ListFlatMapIndexed(l.l, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) FlatMapTo(dest *[]R, f func(T) []R) ISCList[R] {
	r := ListFlatMapTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) FlatMapIndexedTo(dest *[]R, f func(int, T) []R) ISCList[R] {
	r := ListFlatMapIndexedTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) Map(f func(T) R) ISCList[R] {
	r := ListMap(l.l, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) MapIndexed(f func(int, T) R) ISCList[R] {
	r := ListMapIndexed(l.l, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) MapTo(dest *[]R, f func(T) R) ISCList[R] {
	r := ListMapTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) MapIndexedTo(dest *[]R, f func(int, T) R) ISCList[R] {
	r := ListMapIndexedTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCListToMap[T, R]) Reduce(init func(T) R, f func(R, T) R) R {
	return Reduce(l.l, init, f)
}

func (l ISCListToMap[T, R]) ReduceIndexed(init func(int, T) R, f func(int, R, T) R) R {
	return ReduceIndexed(l.l, init, f)
}

type ISCListToGroup[T comparable, K comparable, V comparable] struct {
	l []T
}

func ListToGroupFrom[T comparable, K comparable, V comparable](list ISCList[T]) ISCListToGroup[T, K, V] {
	return ISCListToGroup[T, K, V]{
		l: list.l,
	}
}

func (l ISCListToGroup[T, K, V]) GroupBy(f func(T) K) map[K][]T {
	return GroupBy(l.l, f)
}

func (l ISCListToGroup[T, K, V]) GroupByTransform(f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransform(l.l, f, trans)
}

func (l ISCListToGroup[T, K, V]) GroupByTo(dest *map[K][]T, f func(T) K) map[K][]T {
	return GroupByTo(l.l, dest, f)
}

func (l ISCListToGroup[T, K, V]) GroupByTransformTo(dest *map[K][]V, f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransformTo(l.l, dest, f, trans)
}

type ISCListToPair[K comparable, V comparable] struct {
	l []Pair[K, V]
}

func ListToPairFrom[K comparable, V comparable](list ISCList[Pair[K, V]]) ISCListToPair[K, V] {
	return ISCListToPair[K, V]{
		l: list.l,
	}
}

func ListToPairWithPairs[K comparable, V comparable](list ...Pair[K, V]) ISCListToPair[K, V] {
	return ISCListToPair[K, V]{
		l: list,
	}
}

func (l ISCListToPair[K, V]) ToMap() ISCMap[K, V] {
	m := make(map[K]V)
	for _, item := range l.l {
		m[item.First] = item.Second
	}
	return NewMap(m)
}
