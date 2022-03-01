package isc

// ListEquals 比较两个数组是否相同
func ListEquals[T comparable](leftList []T, rightList []T) bool {
	if leftList == nil && rightList == nil {
		return true
	}
	if leftList == nil || rightList == nil {
		return false
	}
	if len(leftList) != len(rightList) {
		return false
	}
	l := ListDistinct(append(leftList, rightList...))
	return len(l) == len(leftList)
}

// MapEquals 比较两个map是否相同
func MapEquals[K comparable, V comparable](leftMap map[K]V, rightMap map[K]V) bool {
	if leftMap == nil && rightMap == nil {
		return true
	}
	if leftMap == nil || rightMap == nil {
		return false
	}
	if len(leftMap) != len(rightMap) {
		return false
	}

	for key, leftValue := range leftMap {
		rightValue, exist := rightMap[key]
		if !exist || rightValue != leftValue {
			return false
		}
	}
	return true
}
