package isc

func ListPlus[T any](list []T, n []T) []T {
	var t []T
	for _, e := range list {
		t = append(t, e)
	}
	for _, e := range n {
		if !ListContains(t, e) {
			t = append(t, e)
		}
	}
	return t
}

func ListMinus[T any](list []T, n []T) []T {
	var t []T
	for _, e := range list {
		if !ListContains(n, e) {
			t = append(t, e)
		}
	}
	return t
}

func MapPlus[K comparable, V any](m map[K]V, n map[K]V) map[K]V {
	r := make(map[K]V)
	for k, v := range m {
		r[k] = v
	}
	for k, v := range n {
		r[k] = v
	}
	return r
}

func MapMinus[K comparable, V any](m map[K]V, n map[K]V) map[K]V {
	r := make(map[K]V)
	for k, v := range m {
		if _, ok := n[k]; !ok {
			r[k] = v
		}
	}
	return r
}
