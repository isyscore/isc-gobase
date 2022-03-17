package isc

func Int(AStart, AEnd int) []int {
	var ret []int
	for i := AStart; i < AEnd; i++ {
		ret = append(ret, i)
	}
	return ret
}

func Int64(AStart, AEnd int64) []int64 {
	var ret []int64
	for i := AStart; i < AEnd; i++ {
		ret = append(ret, i)
	}
	return ret
}

func IntStep(AStart, AEnd, AStep int) []int {
	var ret []int
	for i := AStart; i < AEnd; i += AStep {
		ret = append(ret, i)
	}
	return ret
}

func Int64Step(AStart, AEnd, AStep int64) []int64 {
	var ret []int64
	for i := AStart; i < AEnd; i += AStep {
		ret = append(ret, i)
	}
	return ret
}

type MapPair[K comparable, V comparable] struct {
	Key   K
	Value V
}

func OrderMapToList[K comparable, V comparable](m OrderMap[K, V]) []MapPair[K, V] {
	var n []MapPair[K, V]
	for _, key := range m.Keys() {
		n = append(n, MapPair[K, V]{key, m.Data[key]})
	}
	return n
}
