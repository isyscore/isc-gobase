package isc

func Find[T any](list []T, f func(T) bool) *T {
	var n *T = nil
	for _, e := range list {
		if f(e) {
			n = &e
			break
		}
	}
	return n
}

func FindLast[T any](list []T, f func(T) bool) *T {
	var n *T = nil
	for i := len(list) - 1; i >= 0; i-- {
		e := list[i]
		if f(e) {
			n = &e
			break
		}
	}
	return n
}

func First[T any](list []T) T {
	return list[0]
}

func Last[T any](list []T) T {
	return list[len(list)-1]
}

func FirstOrNull[T any](list []T) *T {
	var n *T = nil
	if len(list) > 0 {
		n = &list[0]
	}
	return n
}

func LastOrNull[T any](list []T) *T {
	var n *T = nil
	if len(list) > 0 {
		n = &list[len(list)-1]
	}
	return n
}

//IndexOf 判断元素item是否在分片中，示例res := IndexOf[int](list,item),使用时须指明类型
func IndexOf[T comparable](list []T, item T) int {
	var idx = -1
	for i, e := range list {
		if e == item {
			idx = i
			break
		}
	}
	return idx
}

func LastIndexOf[T comparable](list []T, item T) int {
	var idx = -1
	for i := len(list) - 1; i >= 0; i-- {
		e := list[i]
		if e == item {
			idx = i
			break
		}
	}
	return idx
}

func IndexOfCondition[T comparable](list []T, f func(T) bool) int {
	var idx = -1
	for i, e := range list {
		if f(e) {
			idx = i
			break
		}
	}
	return idx
}

func LastIndexOfCondition[T comparable](list []T, f func(T) bool) int {
	var idx = -1
	for i := len(list) - 1; i >= 0; i-- {
		e := list[i]
		if f(e) {
			idx = i
			break
		}
	}
	return idx
}
