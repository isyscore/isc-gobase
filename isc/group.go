package isc

//GroupBy Groups elements of the original collection by the key returned by the given keySelector function applied to each element and returns a map where each group key is associated with a list of corresponding elements.
//The returned map preserves the entry iteration order of the keys produced from the original collection.
func GroupBy[T any, K comparable](list []T, keySelector func(T) K) (destination map[K][]T) {
	dest := make(map[K][]T)
	return GroupByTo(list, &dest, keySelector)
}

//GroupByTransform Groups values returned by the trans function applied to each element of the original collection
//by the key returned by the given keySelector function applied to the element and puts to a map
//each group key associated with a list of corresponding values.
func GroupByTransform[T any, K comparable, V any](list []T, keySelector func(T) K, trans func(T) V) map[K][]V {
	dest := make(map[K][]V)
	return GroupByTransformTo(list, &dest, keySelector, trans)
	//var r = make(map[K][]V)
	//for _, e := range list {
	//	key := keySelector(e)
	//	sl, _ := r[key]
	//	sl = append(sl, trans(e))
	//	r[key] = sl
	//}
	//return r
}

//GroupByTo Groups elements of the original collection by the key returned by the given keySelector function applied
//to each element and puts to the dest map each group key associated with a list of corresponding elements.
//Returns: The dest map‘s val
func GroupByTo[T any, K comparable](list []T, dest *map[K][]T, keySelector func(T) K) (destination map[K][]T) {
	r := make(map[K][]T)
	if *dest == nil {
		return r
	}
	for _, e := range list {
		key := keySelector(e)
		sl, _ := r[key]
		sl = append(sl, e)
		r[key] = sl
		sll, _ := (*dest)[key]
		sll = append(sll, e)
		(*dest)[key] = sll
	}
	return r
}

//GroupByTransformTo  Groups values returned by the trans function applied to each element of the original collection
//by the key returned by the given keySelector function applied to the element and puts to the dest map
//each group key associated with a list of corresponding values.
//Returns: The dest map's val.
func GroupByTransformTo[T any, K comparable, V any](list []T, dest *map[K][]V, keySelector func(T) K, trans func(T) V) map[K][]V {
	r := make(map[K][]V)
	for _, e := range list {
		key := keySelector(e)
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
