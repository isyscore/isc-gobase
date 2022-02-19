package isc

//Take Returns a list containing first n elements.
func Take[T any](list []T, n int) []T {
	var t []T
	if n <= 0 {
		return t
	}
	if n >= len(list) {
		return list
	}
	return SubList(list, 0, n)
}

//TakeLast Returns a list containing last n elements.
func TakeLast[T any](list []T, n int) []T {
	var t []T
	if n <= 0 {
		return t
	}
	if n >= len(list) {
		return list
	}
	return SubList(list, len(list)-n, len(list))
}

//TakeWhile Returns a list containing first elements satisfying the given predicate.
func TakeWhile[T any](list []T, n int, predicate func(T) bool) []T {
	if n <= 0 {
		return nil
	}
	t := make([]T, n)
	for i := 0; i < len(list); i++ {
		if len(t) < n && predicate(list[i]) {
			t = append(t, list[i])
		}
	}
	return t
}

//TakeLastWhile Returns a list containing first elements satisfying the given predicate.
func TakeLastWhile[T any](list []T, n int, predicate func(T) bool) []T {
	if n <= 0 {
		return nil
	}
	t := make([]T, n)
	for i := len(list) - 1; i >= 0; i++ {
		if len(t) < n && predicate(list[i]) {
			t = append(t, list[i])
		}
	}
	return t
}

//Drop Returns a list containing all elements except first [n] elements.
func Drop[T any](list []T, n int) []T {
	if n < 0 {
		return list
	}
	return list[n:]
}

//DropLast Returns a list containing all elements except last n elements
func DropLast[T any](list []T, n int) []T {
	if n < 0 {
		return nil
	}
	if n == 0 {
		return list
	}
	return list[:n]
}

//DropWhile Returns a list containing all elements except first elements that satisfy the given predicate.
func DropWhile[T any](list []T, n int, predicate func(T) bool) []T {
	var t []T
	var dropCount = 0
	for i := 0; i < len(list); i++ {
		if dropCount < n {
			if predicate(list[i]) {
				// dropped!
				dropCount++
			} else {
				t = append(t, list[i])
			}
		} else {
			t = append(t, list[i])
		}
	}
	return t
}

//DropLastWhile Returns a list containing all elements except last elements that satisfy the given predicate.
func DropLastWhile[T any](list []T, n int, predicate func(T) bool) []T {
	var t []T
	var dropCount = 0
	for i := len(list) - 1; i >= 0; i-- {
		if dropCount < n {
			if predicate(list[i]) {
				// dropped!
				dropCount++
			} else {
				t = append(t, list[i])
			}
		} else {
			t = append(t, list[i])
		}
	}
	// reverse
	var r []T
	for i := len(t) - 1; i >= 0; i-- {
		r = append(r, t[i])
	}
	return t
}
