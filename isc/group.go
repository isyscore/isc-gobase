package isc

func GroupBy[T any, K comparable](list []T, f func(T) K) map[K][]T {
	var r = make(map[K][]T)
	for _, e := range list {
		key := f(e)
		sl, _ := r[key]
		sl = append(sl, e)
		r[key] = sl
	}
	return r
}

func GroupByTransform[T any, K comparable, V any](list []T, f func(T) K, trans func(T) V) map[K][]V {
	var r = make(map[K][]V)
	for _, e := range list {
		key := f(e)
		sl, _ := r[key]
		sl = append(sl, trans(e))
		r[key] = sl
	}
	return r
}

func GroupByTo[T any, K comparable](list []T, dest *map[K][]T, f func(T) K) map[K][]T {
	var r = make(map[K][]T)
	for _, e := range list {
		key := f(e)
		sl, _ := r[key]
		sl = append(sl, e)
		r[key] = sl
		sll, _ := (*dest)[key]
		sll = append(sll, e)
		(*dest)[key] = sll
	}
	return r
}

func GroupByTransformTo[T any, K comparable, V any](list []T, dest *map[K][]V, f func(T) K, trans func(T) V) map[K][]V {
	var r = make(map[K][]V)
	for _, e := range list {
		key := f(e)
		value := trans(e)
		sl, _ := r[key]
		sl = append(sl, value)
		r[key] = sl
		sll, _ := (*dest)[key]
		sll = append(sll, value)
		(*dest)[key] = sll
	}
	return r
}
