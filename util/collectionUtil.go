package util

func Contain(dataList []interface{}, data interface{}) bool {
	for _, item := range dataList {
		if item == data {
			return true
		}
	}
	return false
}
