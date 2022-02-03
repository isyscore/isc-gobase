package isc

func SubList[T any](list []T, fromIndex int, toIndex int) []T {
	var n []T
	for i := 0; i < len(list); i++ {
		if i >= fromIndex && i < toIndex {
			n = append(n, list[i])
		}
	}
	return n
}

func Slice[T any](list []T, r IntRange) []T {
	var n []T
	for i := 0; i < len(list); i++ {
		if i >= r.Start && i < r.End {
			n = append(n, list[i])
		}
	}
	return n
}

func SliceBy[T comparable](list []T, r []int) []T {
	var n []T
	for i := 0; i < len(list); i++ {
		if ListContains(r, i) {
			n = append(n, list[i])
		}
	}
	return n
}
