package isc

func ListJoinToStringFull[T any](list []T, sep string, prefix string, postfix string, f func(T) string) string {
	buffer := prefix
	var count = 0
	for _, e := range list {
		count++
		if count > 1 {
			buffer += sep
		}
		buffer += f(e)
	}
	buffer += postfix
	return buffer
}

func ListJoinToString[T any](list []T, f func(T) string) string {
	return ListJoinToStringFull(list, ",", "", "", f)
}

func MapJoinToStringFull[K comparable, V any](m map[K]V, sep string, prefix string, postfix string, f func(K, V) string) string {
	buffer := prefix
	var count = 0
	for k, v := range m {
		count++
		if count > 1 {
			buffer += sep
		}
		buffer += f(k, v)
	}
	buffer += postfix
	return buffer
}

func MapJoinToString[K comparable, V any](m map[K]V, f func(K, V) string) string {
	return MapJoinToStringFull(m, ",", "", "", f)
}
