package isc

//Associate Returns a Map containing key-value map provided by transform function applied to elements of the given collection.
//If any of two map would have the same key the last one gets added to the map.
//The returned map preserves the entry iteration order of the original collection.
func Associate[T any, K comparable, V any](list []T, transform func(T) Pair[K, V]) map[K]V {
	r := make(map[K]V)
	return AssociateTo(list, &r, transform)
}

//AssociateTo Populates and returns the destination map with key-value pairs provided by transform function applied to each element of the given collection.
//If any of two pairs would have the same key the last one gets added to the map.
func AssociateTo[T any, K comparable, V any](list []T, destination *map[K]V, transform func(T) Pair[K, V]) map[K]V {
	for _, e := range list {
		item := transform(e)
		(*destination)[item.First] = item.Second
	}
	return *destination
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
func AssociateByAndValue[T any, V any, K comparable](list []T, keySelector func(T) K, valueTransform func(T) V) map[K]V {
	r := make(map[K]V)
	for _, e := range list {
		r[keySelector(e)] = valueTransform(e)
	}
	return r
}

// AssociateByTo Populates and returns the destination mutable map with key-value pairs, where key is provided by the keySelector function applied to each element of the given collection and value is the element itself.
//If any two elements would have the same key returned by keySelector the last one gets added to the map
func AssociateByTo[T any, K comparable](list []T, destination *map[K]T, keySelector func(T) K) map[K]T {
	for _, e := range list {
		(*destination)[keySelector(e)] = e
	}
	return *destination
}

// AssociateByAndValueTo Populates and returns the destination mutable map with key-value pairs, where key is provided by the keySelector function applied to each element of the given collection and value is the element itself.
//If any two elements would have the same key returned by keySelector the last one gets added to the map
func AssociateByAndValueTo[T, V any, K comparable](list []T, destination *map[K]V, keySelector func(T) K, valueTransform func(T) V) map[K]V {
	for _, e := range list {
		(*destination)[keySelector(e)] = valueTransform(e)
	}
	return *destination
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
func AssociateWithTo[T comparable, V any](list []T, destination *map[T]V, valueSelector func(T) V) map[T]V {
	for _, e := range list {
		(*destination)[e] = valueSelector(e)
	}
	return *destination
}
