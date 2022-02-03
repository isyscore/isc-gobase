package isc

func ListFlatMap[T any, R any](list []T, f func(T) []R) []R {
	var r []R
	for _, e := range list {
		rlist := f(e)
		for _, rl := range rlist {
			r = append(r, rl)
		}
	}
	return r
}

func ListFlatMapIndexed[T any, R any](list []T, f func(int, T) []R) []R {
	var r []R
	for i, e := range list {
		rlist := f(i, e)
		for _, rl := range rlist {
			r = append(r, rl)
		}
	}
	return r
}

func ListFlattern[T any](list [][]T) []T {
	var r []T
	for _, e := range list {
		rlist := e
		for _, rl := range rlist {
			r = append(r, rl)
		}
	}
	return r
}

func ListFlatMapTo[T any, R any](list []T, dest *[]R, f func(T) []R) []R {
	var r []R
	for _, e := range list {
		rlist := f(e)
		for _, rl := range rlist {
			*dest = append(*dest, rl)
			r = append(r, rl)
		}
	}
	return r
}

func ListFlatMapIndexedTo[T any, R any](list []T, dest *[]R, f func(int, T) []R) []R {
	var r []R
	for i, e := range list {
		rlist := f(i, e)
		for _, rl := range rlist {
			*dest = append(*dest, rl)
			r = append(r, rl)
		}
	}
	return r
}

func MapFlatMap[K comparable, V any, R any](m map[K]V, f func(K, V) []R) []R {
	var r []R
	for k, v := range m {
		rlist := f(k, v)
		for _, rl := range rlist {
			r = append(r, rl)
		}
	}
	return r
}

func MapFlatMapTo[K comparable, V any, R any](m map[K]V, dest *[]R, f func(K, V) []R) []R {
	var r []R
	for k, v := range m {
		rlist := f(k, v)
		for _, rl := range rlist {
			*dest = append(*dest, rl)
			r = append(r, rl)
		}
	}
	return r
}

