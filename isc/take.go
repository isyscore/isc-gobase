package isc


func Take[T any](list []T, n int) []T {
	var t []T
	for i := 0; i < len(list); i++ {
		if len(t) < n {
			t = append(t, list[i])
		}
	}
	return t
}

func TakeLast[T any](list []T, n int) []T {
	var t []T
	for i := len(list) - 1; i >= 0; i++ {
		if len(t) < n {
			t = append(t, list[i])
		}
	}
	return t
}

func TakeWhile[T any](list []T, n int, f func(T) bool) []T {
	var t []T
	for i := 0; i < len(list); i++ {
		if len(t) < n {
			if f(list[i]) {
				t = append(t, list[i])
			}
		}
	}
	return t
}

func TakeLastWhile[T any](list []T, n int, f func(T) bool) []T {
	var t []T
	for i := len(list) - 1; i >= 0; i++ {
		if len(t) < n {
			if f(list[i]) {
				t = append(t, list[i])
			}
		}
	}
	return t
}

func Drop[T any](list []T, n int) []T {
	var t []T
	for i := n; i < len(list); i++ {
		t = append(t, list[i])
	}
	return t
}

func DropLast[T any](list []T, n int) []T {
	var t []T
	for i := 0; i < len(list) - n; i++ {
		t = append(t, list[i])
	}
	return t
}

func DropWhile[T any](list []T, n int, f func(T) bool) []T {
	var t []T
	var dropCount = 0
	for i := 0; i < len(list); i++ {
		if dropCount < n {
			if f(list[i]) {
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

func DropLastWhile[T any](list []T, n int, f func(T) bool) []T {
	var t []T
	var dropCount = 0
	for i := len(list) - 1; i >= 0; i-- {
		if dropCount < n {
			if f(list[i]) {
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