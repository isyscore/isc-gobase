package isc

type LinkedHashMap[K comparable, V comparable] OrderMap[K, V]

type OrderMap[K comparable, V comparable] struct {
	Data    map[K]V
	KeyList []K
}

func NewOrderMap[K comparable, V comparable]() OrderMap[K, V] {
	return OrderMap[K, V]{
		Data:    make(map[K]V),
		KeyList: []K{},
	}
}

func (m OrderMap[K, V]) Size() int {
	return len(m.KeyList)
}

func (m *OrderMap[K, V]) Put(k K, v V) {
	if !ListContains(m.KeyList, k) {
		m.KeyList = append(m.KeyList, k)
	}
	m.Data[k] = v
}

func (m *OrderMap[K, V]) PutPair(item Pair[K, V]) {
	if !ListContains(m.KeyList, item.First) {
		m.KeyList = append(m.KeyList, item.First)
	}
	m.Data[item.First] = item.Second
}

func (m *OrderMap[K, V]) PutPairs(item ...Pair[K, V]) {
	for _, pair := range item {
		m.PutPair(pair)
	}
}

func (m OrderMap[K, V]) Get(k K) V {
	return m.Data[k]
}

func (m OrderMap[K, V]) GetOrDef(k K, def V) V {
	if v, ok := m.Data[k]; ok {
		return v
	} else {
		return def
	}
}

func (m *OrderMap[K, V]) Delete(k K) {
	idx := IndexOf(m.KeyList, k)
	if idx != -1 {
		m.KeyList = append(m.KeyList[:idx], m.KeyList[idx+1:]...)
		delete(m.Data, k)
	}
}

func (m *OrderMap[K, V]) Clear() {
	m.KeyList = []K{}
	m.Data = make(map[K]V)
}

func (m OrderMap[K, V]) Keys() []K {
	return m.KeyList
}

func (m OrderMap[K, V]) GetKey(index int) K {
	return m.KeyList[index]
}

func (m OrderMap[K, V]) GetValue(index int) V {
	key := m.KeyList[index]
	return m.Data[key]
}

func (m OrderMap[K, V]) ForEach(f func(K, V)) {
	for _, key := range m.KeyList {
		f(key, m.Data[key])
	}
}

func (m OrderMap[K, V]) ForEachIndexed(f func(int, K, V)) {
	for idx, key := range m.KeyList {
		f(idx, key, m.Data[key])
	}
}

func (m OrderMap[K, V]) Filter(f func(K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if f(key, m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterIndexed(f func(int, K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if f(idx, key, m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterNot(f func(K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if !f(key, m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterNotIndexed(f func(int, K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if !f(idx, key, m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterKeys(f func(K) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if f(key) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterKeysIndexed(f func(int, K) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if f(idx, key) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterValues(f func(V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if f(m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterValuesIndexed(f func(int, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if f(idx, m.Data[key]) {
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterTo(dest *OrderMap[K, V], f func(K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if f(key, m.Data[key]) {
			dest.Put(key, m.Data[key])
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterIndexedTo(dest *OrderMap[K, V], f func(int, K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if f(idx, key, m.Data[key]) {
			dest.Put(key, m.Data[key])
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterNotTo(dest *OrderMap[K, V], f func(K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if !f(key, m.Data[key]) {
			dest.Put(key, m.Data[key])
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) FilterNotIndexedTo(dest *OrderMap[K, V], f func(int, K, V) bool) OrderMap[K, V] {
	result := NewOrderMap[K, V]()
	for idx, key := range m.KeyList {
		if !f(idx, key, m.Data[key]) {
			dest.Put(key, m.Data[key])
			result.Put(key, m.Data[key])
		}
	}
	return result
}

func (m OrderMap[K, V]) Contains(k K, v V) bool {
	return MapContains(m.Data, k, v)
}

func (m OrderMap[K, V]) ContainsKey(k K) bool {
	return MapContainsKey(m.Data, k)
}

func (m OrderMap[K, V]) ContainsValue(v V) bool {
	return MapContainsValue(m.Data, v)
}

func (m OrderMap[K, V]) JoinToString(f func(K, V) string) string {
	return m.JoinToStringFull(",", "", "", f)
}

func (m OrderMap[K, V]) JoinToStringFull(sep string, prefix string, postfix string, f func(K, V) string) string {
	buffer := prefix
	var count = 0
	for _, key := range m.KeyList {
		count++
		if count > 1 {
			buffer += sep
		}
		buffer += f(key, m.Data[key])
	}
	buffer += postfix
	return buffer
}

func (m OrderMap[K, V]) All(f func(K, V) bool) bool {
	for _, key := range m.KeyList {
		if !f(key, m.Data[key]) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) Any(f func(K, V) bool) bool {
	for _, key := range m.KeyList {
		if f(key, m.Data[key]) {
			return true
		}
	}
	return false
}

func (m OrderMap[K, V]) None(f func(K, V) bool) bool {
	for _, key := range m.KeyList {
		if f(key, m.Data[key]) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) Count(f func(K, V) bool) int {
	num := 0
	for _, key := range m.KeyList {
		if f(key, m.Data[key]) {
			num++
		}
	}
	return num
}

func (m OrderMap[K, V]) AllKey(f func(K) bool) bool {
	for _, key := range m.KeyList {
		if !f(key) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) AnyKey(f func(K) bool) bool {
	for _, key := range m.KeyList {
		if f(key) {
			return true
		}
	}
	return false
}

func (m OrderMap[K, V]) NoneKey(f func(K) bool) bool {
	for _, key := range m.KeyList {
		if f(key) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) CountKey(f func(K) bool) int {
	num := 0
	for _, key := range m.KeyList {
		if f(key) {
			num++
		}
	}
	return num
}

func (m OrderMap[K, V]) AllValue(f func(V) bool) bool {
	for _, key := range m.KeyList {
		if !f(m.Data[key]) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) AnyValue(f func(V) bool) bool {
	for _, key := range m.KeyList {
		if f(m.Data[key]) {
			return true
		}
	}
	return false
}

func (m OrderMap[K, V]) NoneValue(f func(V) bool) bool {
	for _, key := range m.KeyList {
		if f(m.Data[key]) {
			return false
		}
	}
	return true
}

func (m OrderMap[K, V]) CountValue(f func(V) bool) int {
	num := 0
	for _, key := range m.KeyList {
		if f(m.Data[key]) {
			num++
		}
	}
	return num
}

func (m OrderMap[K, V]) ToList() []Pair[K, V] {
	var n []Pair[K, V]
	for _, key := range m.KeyList {
		n = append(n, NewPair(key, m.Data[key]))
	}
	return n
}

func (m OrderMap[K, V]) Plus(n OrderMap[K, V]) OrderMap[K, V] {
	r := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		r.Put(key, m.Data[key])
	}
	n.ForEach(func(k K, v V) {
		r.Put(k, v)
	})
	return r
}

func (m OrderMap[K, V]) Minus(n OrderMap[K, V]) OrderMap[K, V] {
	r := NewOrderMap[K, V]()
	for _, key := range m.KeyList {
		if _, ok := n.Data[key]; !ok {
			r.Put(key, m.Data[key])
		}
	}
	return r
}
