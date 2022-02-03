package isc

/// functions for list

func ListMap[T any, R any](list []T, f func(T) R) []R {
	var n []R
	for _, e := range list {
		n = append(n, f(e))
	}
	return n
}

func ListMapNotNull[T any, R any](list []*T, f func(T) R) []R {
	var n []R
	for _, e := range list {
		if e != nil {
			n = append(n, f(*e))
		}
	}
	return n
}

func ListMapIndexed[T any, R any](list []T, f func(int, T) R) []R {
	var n []R
	for i, e := range list {
		n = append(n, f(i, e))
	}
	return n
}

func ListMapIndexedNotNull[T any, R any](list []*T, f func(int, T) R) []R {
	var n []R
	for i, e := range list {
		if e != nil {
			n = append(n, f(i, *e))
		}
	}
	return n
}

func ListMapTo[T any, R any](list []T, dest *[]R, f func(T) R) []R {
	var n []R
	for _, e := range list {
		item := f(e)
		*dest = append(*dest, item)
		n = append(n, item)
	}
	return n
}

func ListMapIndexedTo[T any, R any](list []T, dest *[]R, f func(int, T) R) []R {
	var n []R
	for i, e := range list {
		item := f(i, e)
		*dest = append(*dest, item)
		n = append(n, item)
	}
	return n
}

func ListMapNotNullTo[T any, R any](list []*T, dest *[]R, f func(T) R) []R {
	var n []R
	for _, e := range list {
		if e != nil {
			item := f(*e)
			*dest = append(*dest, item)
			n = append(n, item)
		}
	}
	return n
}

func ListMapIndexedNotNullTo[T any, R any](list []*T, dest *[]R, f func(int, T) R) []R {
	var n []R
	for i, e := range list {
		if e != nil {
			item := f(i, *e)
			*dest = append(*dest, item)
			n = append(n, item)
		}
	}
	return n
}

/// functions for map

func MapMap[K comparable, V any, R any](m map[K]V, f func(K, V) R) []R {
	var n []R
	for k, v := range m {
		n = append(n, f(k, v))
	}
	return n
}

func MapMapNotNull[K comparable, V any, R any](m map[K]*V, f func(K, V) R) []R {
	var n []R
	for k, v := range m {
		if v != nil {
			n = append(n, f(k, *v))
		}
	}
	return n
}

func MapMapTo[K comparable, V any, R any](m map[K]V, dest *[]R, f func(K, V) R) []R {
	var n []R
	for k, v := range m {
		item := f(k, v)
		*dest = append(*dest, item)
		n = append(n, item)
	}
	return n
}

func MapMapNotNullTo[K comparable, V any, R any](m map[K]*V, dest *[]R, f func(K, V) R) []R {
	var n []R
	for k, v := range m {
		if v != nil {
			item := f(k, *v)
			*dest = append(*dest, item)
			n = append(n, item)
		}
	}
	return n
}
