package isc

type ISCListToMap[T any, R any] struct {
	ISCList[T]
}

func ListToMapFrom[T any, R any](list ISCList[T]) ISCListToMap[T, R] {
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

type ISCListToSlice[T any, R comparable] struct {
	ISCList[T]
}

func ListToSliceFrom[T any, R comparable](list ISCList[T]) ISCListToSlice[T, R] {
	return ISCListToSlice[T, R]{
		list,
	}
}

func (l ISCListToSlice[T, R]) SliceContains(predicate func(T) R, key R) bool {
	return SliceContains(l.ISCList, predicate, key)
}

func (l ISCListToSlice[T, R]) SliceTo(valueTransform func(T) R) ISCMap[R, T] {
	return SliceTo(l.ISCList, valueTransform)
}

type ISCListToTriple[T comparable, K comparable, V any] struct {
	ISCList[T]
}

func ListToTripleFrom[T comparable, K comparable, V any](list ISCList[T]) ISCListToTriple[T, K, V] {
	return ISCListToTriple[T, K, V]{
		list,
	}
}

func (l ISCListToTriple[T, K, V]) GroupBy(f func(T) K) map[K][]T {
	return GroupBy(l.ISCList, f)
}

func (l ISCListToTriple[T, K, V]) GroupByTransform(f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransform(l.ISCList, f, trans)
}

func (l ISCListToTriple[T, K, V]) GroupByTo(dest *map[K][]T, f func(T) K) map[K][]T {
	return GroupByTo(l.ISCList, dest, f)
}

func (l ISCListToTriple[T, K, V]) GroupByTransformTo(dest *map[K][]V, f func(T) K, trans func(T) V) map[K][]V {
	return GroupByTransformTo(l.ISCList, dest, f, trans)
}

func (l ISCListToTriple[T, K, V]) Associate(transform func(T) Pair[K, V]) ISCMap[K, V] {
	return Associate(l.ISCList, transform)
}

func (l ISCListToTriple[T, K, V]) AssociateTo(destination *map[K]V, transform func(T) Pair[K, V]) ISCMap[K, V] {
	return AssociateTo(l.ISCList, destination, transform)
}

func (l ISCListToTriple[T, K, V]) AssociateBy(keySelector func(T) K) ISCMap[K, T] {
	return AssociateBy(l.ISCList, keySelector)
}

func (l ISCListToTriple[T, K, V]) AssociateByAndValue(keySelector func(T) K, valueTransform func(T) V) ISCMap[K, V] {
	return AssociateByAndValue(l.ISCList, keySelector, valueTransform)
}

func (l ISCListToTriple[T, K, V]) AssociateByTo(destination *map[K]T, keySelector func(T) K) ISCMap[K, T] {
	return AssociateByTo(l.ISCList, destination, keySelector)
}

func (l ISCListToTriple[T, K, V]) AssociateByAndValueTo(destination *map[K]V, keySelector func(T) K, valueTransform func(T) V) ISCMap[K, V] {
	return AssociateByAndValueTo(l.ISCList, destination, keySelector, valueTransform)
}

func (l ISCListToTriple[T, K, V]) AssociateWith(valueSelector func(T) V) ISCMap[T, V] {
	return AssociateWith(l.ISCList, valueSelector)
}

func (l ISCListToTriple[T, K, V]) AssociateWithTo(destination *map[T]V, valueSelector func(T) V) ISCMap[T, V] {
	return AssociateWithTo(l.ISCList, destination, valueSelector)
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
