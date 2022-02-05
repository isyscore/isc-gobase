package isc

type ISCList[T comparable] struct {
	l []T
}

func NewList[T comparable](list []T) ISCList[T] {
	return ISCList[T]{l: list}
}

func NewListWithItems[T comparable](items ...T) ISCList[T] {
	return ISCList[T]{l: items}
}

func (l ISCList[T]) ToArray() []T {
	return l.l
}

func (l ISCList[T]) IsEmpty() bool {
	return len(l.l) == 0
}

func (l ISCList[T]) Filter(f func(T) bool) ISCList[T] {
	r := ListFilter(l.l, f)
	return NewList(r)
}

func (l ISCList[T]) FilterNot(f func(T) bool) ISCList[T] {
	r := ListFilterNot(l.l, f)
	return NewList(r)
}

func (l ISCList[T]) FilterIndexed(f func(int, T) bool) ISCList[T] {
	r := ListFilterIndexed(l.l, f)
	return NewList(r)
}

func (l ISCList[T]) FilterNotIndexed(f func(int, T) bool) ISCList[T] {
	r := ListFilterNotIndexed(l.l, f)
	return NewList(r)
}

func (l ISCList[T]) FilterTo(dest *[]T, f func(T) bool) ISCList[T] {
	r := ListFilterTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCList[T]) FilterNotTo(dest *[]T, f func(T) bool) ISCList[T] {
	r := ListFilterNotTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCList[T]) FilterIndexedTo(dest *[]T, f func(int, T) bool) ISCList[T] {
	r := ListFilterIndexedTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCList[T]) FilterNotIndexedTo(dest *[]T, f func(int, T) bool) ISCList[T] {
	r := ListFilterNotIndexedTo(l.l, dest, f)
	return NewList(r)
}

func (l ISCList[T]) Contains(item T) bool {
	return ListContains(l.l, item)
}

func (l ISCList[T]) Find(f func(T) bool) *T {
	return Find(l.l, f)
}

func (l ISCList[T]) FindLast(f func(T) bool) *T {
	return FindLast(l.l, f)
}

func (l ISCList[T]) First() T {
	return First(l.l)
}

func (l ISCList[T]) Last() T {
	return Last(l.l)
}

func (l ISCList[T]) FirstOrNull() *T {
	return FirstOrNull(l.l)
}

func (l ISCList[T]) LastOrNull() *T {
	return LastOrNull(l.l)
}

func (l ISCList[T]) IndexOf(item T) int {
	return IndexOf(l.l, item)
}

func (l ISCList[T]) LastIndexOf(item T) int {
	return LastIndexOf(l.l, item)
}

func (l ISCList[T]) IndexOfCondition(f func(T) bool) int {
	return IndexOfCondition(l.l, f)
}

func (l ISCList[T]) LastIndexOfCondition(f func(T) bool) int {
	return LastIndexOfCondition(l.l, f)
}

func (l ISCList[T]) JoinToString(f func(T) string) string {
	return ListJoinToString(l.l, f)
}

func (l ISCList[T]) JoinToStringFull(sep string, prefix string, postfix string, f func(T) string) string {
	return ListJoinToStringFull(l.l, sep, prefix, postfix, f)
}

func (l ISCList[T]) All(f func(T) bool) bool {
	return ListAll(l.l, f)
}

func (l ISCList[T]) Any(f func(T) bool) bool {
	return ListAny(l.l, f)
}

func (l ISCList[T]) None(f func(T) bool) bool {
	return ListNone(l.l, f)
}

func (l ISCList[T]) Count(f func(T) bool) int {
	return ListCount(l.l, f)
}

func (l ISCList[T]) SubList(fromIndex int, toIndex int) ISCList[T] {
	r := SubList(l.l, fromIndex, toIndex)
	return NewList(r)
}

func (l ISCList[T]) Slice(r IntRange) ISCList[T] {
	rr := Slice(l.l, r)
	return NewList(rr)
}

func (l ISCList[T]) SliceBy(r []int) ISCList[T] {
	rr := SliceBy(l.l, r)
	return NewList(rr)
}

func (l ISCList[T]) Take(n int) ISCList[T] {
	r := Take(l.l, n)
	return NewList(r)
}

func (l ISCList[T]) TakeLast(n int) ISCList[T] {
	r := TakeLast(l.l, n)
	return NewList(r)
}

func (l ISCList[T]) TakeWhile(n int, f func(T) bool) ISCList[T] {
	r := TakeWhile(l.l, n, f)
	return NewList(r)
}

func (l ISCList[T]) TakeLastWhile(n int, f func(T) bool) ISCList[T] {
	r := TakeLastWhile(l.l, n, f)
	return NewList(r)
}

func (l ISCList[T]) Drop(n int) ISCList[T] {
	r := Drop(l.l, n)
	return NewList(r)
}

func (l ISCList[T]) DropLast(n int) ISCList[T] {
	r := DropLast(l.l, n)
	return NewList(r)
}

func (l ISCList[T]) DropWhile(n int, f func(T) bool) ISCList[T] {
	r := DropWhile(l.l, n, f)
	return NewList(r)
}

func (l ISCList[T]) DropLastWhile(n int, f func(T) bool) ISCList[T] {
	r := DropLastWhile(l.l, n, f)
	return NewList(r)
}
