package isc

type ISCList[T any] struct {
	l []T
}

func NewList[T any](list []T) *ISCList[T] {
	return &ISCList[T]{l: list}
}

func (l *ISCList[T]) Filter(f func(T) bool) *ISCList[T] {
	r := ListFilter(l.l, f)
	return NewList(r)
}
