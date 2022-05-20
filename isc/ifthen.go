package isc

func IfThen[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
