package isc

type ISCList[T any] []T

func NewList[T any]() ISCList[T] {
	return []T{}
}

func NewListWithList[T any](list []T) ISCList[T] {
	return list
}

func NewListWithItems[T any](items ...T) ISCList[T] {
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

func (l ISCList[T]) Size() int {
	return len(l)
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

func (l ISCList[T]) Partition(partition int) [][]T {
	return Partition(l, partition)
}

func (l ISCList[T]) PartitionWithCal(f func(int) int) [][]T {
	return PartitionWithCal(l, f)
}

func (l ISCList[T]) Plus(n []T) ISCList[T] {
	return ListPlus(l, n)
}

func (l ISCList[T]) Minus(n []T) ISCList[T] {
	return ListMinus(l, n)
}

func (l ISCList[T]) Equals(n ISCList[T]) bool {
	return ListEquals(l, n)
}

func ListToSet[T comparable](list ISCList[T]) ISCSet[T] {
	res := ISCSet[T]{}
	for _, v := range list {
		_ = res.Add(v)
	}
	return res
}
