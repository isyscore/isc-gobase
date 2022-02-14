package isc

type ISCList[T comparable] []T

func NewList[T comparable]() ISCList[T] {
	return []T{}
}

func NewListWithList[T comparable](list []T) ISCList[T] {
	return list
}

func NewListWithItems[T comparable](items ...T) ISCList[T] {
	return items
}

func (l *ISCList[T]) Add(item T) int {
	idx := len(*l)
	*l = append(*l, item)
	return idx
}

func (l *ISCList[T]) AddAll(item ...T) {
	*l = append(*l, item...)
}

func (l *ISCList[T]) Insert(index int, item T) int {
	*l = append((*l)[:index], append([]T{item}, (*l)[index:]...)...)
	return index
}

func (l *ISCList[T]) Delete(index int) T {
	item := (*l)[index]
	*l = append((*l)[:index], (*l)[index+1:]...)
	return item
}

func (l *ISCList[T]) Clear() {
	*l = []T{}
}

func (l ISCList[T]) IsEmpty() bool {
	return len(l) == 0
}

func (l ISCList[T]) ForEach(f func(T)) {
	for _, item := range l {
		f(item)
	}
}

func (l ISCList[T]) ForEachIndexed(f func(int, T)) {
	for idx, item := range l {
		f(idx, item)
	}
}

func (l ISCList[T]) Distinct() ISCList[T] {
	return ListDistinct(l)
}

func (l ISCList[T]) Filter(f func(T) bool) ISCList[T] {
	return ListFilter(l, f)
}

func (l ISCList[T]) FilterNot(f func(T) bool) ISCList[T] {
	return ListFilterNot(l, f)
}

func (l ISCList[T]) FilterIndexed(f func(int, T) bool) ISCList[T] {
	return ListFilterIndexed(l, f)
}

func (l ISCList[T]) FilterNotIndexed(f func(int, T) bool) ISCList[T] {
	return ListFilterNotIndexed(l, f)
}

func (l ISCList[T]) FilterTo(dest *[]T, f func(T) bool) ISCList[T] {
	return ListFilterTo(l, dest, f)
}

func (l ISCList[T]) FilterNotTo(dest *[]T, f func(T) bool) ISCList[T] {
	return ListFilterNotTo(l, dest, f)
}

func (l ISCList[T]) FilterIndexedTo(dest *[]T, f func(int, T) bool) ISCList[T] {
	return ListFilterIndexedTo(l, dest, f)
}

func (l ISCList[T]) FilterNotIndexedTo(dest *[]T, f func(int, T) bool) ISCList[T] {
	return ListFilterNotIndexedTo(l, dest, f)
}

func (l ISCList[T]) Contains(item T) bool {
	return ListContains(l, item)
}

func (l ISCList[T]) Find(f func(T) bool) *T {
	return Find(l, f)
}

func (l ISCList[T]) FindLast(f func(T) bool) *T {
	return FindLast(l, f)
}

func (l ISCList[T]) First() T {
	return First(l)
}

func (l ISCList[T]) Last() T {
	return Last(l)
}

func (l ISCList[T]) FirstOrNull() *T {
	return FirstOrNull(l)
}

func (l ISCList[T]) LastOrNull() *T {
	return LastOrNull(l)
}

func (l ISCList[T]) IndexOf(item T) int {
	return IndexOf(l, item)
}

func (l ISCList[T]) LastIndexOf(item T) int {
	return LastIndexOf(l, item)
}

func (l ISCList[T]) IndexOfCondition(f func(T) bool) int {
	return IndexOfCondition(l, f)
}

func (l ISCList[T]) LastIndexOfCondition(f func(T) bool) int {
	return LastIndexOfCondition(l, f)
}

func (l ISCList[T]) JoinToString(f func(T) string) string {
	return ListJoinToString(l, f)
}

func (l ISCList[T]) JoinToStringFull(sep string, prefix string, postfix string, f func(T) string) string {
	return ListJoinToStringFull(l, sep, prefix, postfix, f)
}

func (l ISCList[T]) All(f func(T) bool) bool {
	return ListAll(l, f)
}

func (l ISCList[T]) Any(f func(T) bool) bool {
	return ListAny(l, f)
}

func (l ISCList[T]) None(f func(T) bool) bool {
	return ListNone(l, f)
}

func (l ISCList[T]) Count(f func(T) bool) int {
	return ListCount(l, f)
}

func (l ISCList[T]) SubList(fromIndex int, toIndex int) ISCList[T] {
	r := SubList(l, fromIndex, toIndex)
	return NewListWithList(r)
}

func (l ISCList[T]) Slice(r IntRange) ISCList[T] {
	rr := Slice(l, r)
	return NewListWithList(rr)
}

func (l ISCList[T]) SliceBy(r []int) ISCList[T] {
	rr := SliceBy(l, r)
	return NewListWithList(rr)
}

func (l ISCList[T]) Take(n int) ISCList[T] {
	r := Take(l, n)
	return NewListWithList(r)
}

func (l ISCList[T]) TakeLast(n int) ISCList[T] {
	r := TakeLast(l, n)
	return NewListWithList(r)
}

func (l ISCList[T]) TakeWhile(n int, f func(T) bool) ISCList[T] {
	r := TakeWhile(l, n, f)
	return NewListWithList(r)
}

func (l ISCList[T]) TakeLastWhile(n int, f func(T) bool) ISCList[T] {
	r := TakeLastWhile(l, n, f)
	return NewListWithList(r)
}

func (l ISCList[T]) Drop(n int) ISCList[T] {
	r := Drop(l, n)
	return NewListWithList(r)
}

func (l ISCList[T]) DropLast(n int) ISCList[T] {
	r := DropLast(l, n)
	return NewListWithList(r)
}

func (l ISCList[T]) DropWhile(n int, f func(T) bool) ISCList[T] {
	r := DropWhile(l, n, f)
	return NewListWithList(r)
}

func (l ISCList[T]) DropLastWhile(n int, f func(T) bool) ISCList[T] {
	r := DropLastWhile(l, n, f)
	return NewListWithList(r)
}

//AssociateBy Returns a Map containing the elements from the given collection indexed by the key returned from keySelector function applied to each element.
//If any two elements would have the same key returned by keySelector the last one gets added to the map.
//The returned map preserves the entry iteration order of the original collection.
func AssociateBy[T any, K comparable](list []T, keySelector func(T) K) map[K]T {
	r := make(map[K]T)
	for _, e := range list {
		r[keySelector(e)] = e
	}
	return r
}

//AssociateByAndValue Returns a Map containing the values provided by valueTransform and indexed by keySelector functions applied to elements of the given collection.
//If any two elements would have the same key returned by keySelector the last one gets added to the map.
//The returned map preserves the entry iteration order of the original collection.
func AssociateByAndValue[T, V any, K comparable](list []T, keySelector func(T) K, valueTransform func(T) V) map[K]V {
	r := make(map[K]V)
	for _, e := range list {
		r[keySelector(e)] = valueTransform(e)
	}
	return r
}

// AssociateByTo Populates and returns the destination mutable map with key-value pairs, where key is provided by the keySelector function applied to each element of the given collection and value is the element itself.
//If any two elements would have the same key returned by keySelector the last one gets added to the map
func AssociateByTo[T, K comparable](list []T, keySelector func(T) K, destination map[K]T) map[K]T {
	for _, e := range list {
		destination[keySelector(e)] = e
	}
	return destination
}

// AssociateByAndValueTo Populates and returns the destination mutable map with key-value pairs, where key is provided by the keySelector function applied to each element of the given collection and value is the element itself.
//If any two elements would have the same key returned by keySelector the last one gets added to the map
func AssociateByAndValueTo[T, V any, K comparable](list []T, keySelector func(T) K, valueTransform func(T) V, destination map[K]V) map[K]V {
	for _, e := range list {
		destination[keySelector(e)] = valueTransform(e)
	}
	return destination
}

//AssociateWith Returns a Map where keys are elements from the given collection and values are produced by the valueSelector function applied to each element.
//If any two elements are equal, the last one gets added to the map.
//The returned map preserves the entry iteration order of the original collection.
func AssociateWith[T comparable, V any](list []T, valueSelector func(T) V) map[T]V {
	destination := make(map[T]V)
	for _, e := range list {
		destination[e] = valueSelector(e)
	}
	return destination
}

//AssociateWithTo Populates and returns the destination mutable map with key-value pairs for each element of the given collection, where key is the element itself and value is provided by the valueSelector function applied to that key.
//If any two elements are equal, the last one overwrites the former value in the map.
func AssociateWithTo[T comparable, V any](list []T, destination *map[T]V, valueSelector func(T) V) *map[T]V {
	for _, e := range list {
		(*destination)[e] = valueSelector(e)
	}
	return destination
}
