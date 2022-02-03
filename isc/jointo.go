package isc

func JoinToStringFull[T any](list []T, sep string, prefix string, postfix string, f func(T) string) string {
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

func JoinToString[T any](list []T, f func(T) string) string {
	return JoinToStringFull(list, ",", "", "", f)
}
