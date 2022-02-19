package isc

/// functions for list

func Reduce[S any, T any](list []T, init func(T) S, f func(S, T) S) S {
	accumulator := init(list[0])
	for _, e := range list[1:] {
		accumulator = f(accumulator, e)
	}
	return accumulator
}

func ReduceIndexed[S any, T any](list []T, init func(int, T) S, f func(int, S, T) S) S {
	accumulator := init(0, list[0])
	for i, e := range list[1:] {
		accumulator = f(i, accumulator, e)
	}
	return accumulator
}
