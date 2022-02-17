package util

// EqualMap 比较两个map是否相同
func EqualMap(leftMap map[string]interface{}, rightMap map[string]interface{}) bool {
	if leftMap == nil && rightMap == nil {
		return true
	}

	if leftMap == nil || rightMap == nil {
		return false
	}

	for key := range leftMap {
		rightValue, exist := rightMap[key]
		if !exist || rightValue != leftMap[key] {
			return false
		}
	}
	return true
}
