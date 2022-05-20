package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestOrderMap(t *testing.T) {
	om := isc.NewOrderMap[string, string]()
	om.Put("a", "1")
	om.Put("b", "2")
	om.Put("c", "3")
	t.Logf("size : %d", om.Size())

	for i := 0; i < om.Size(); i++ {
		t.Logf("key : %s, value : %s", om.GetKey(i), om.GetValue(i))
	}

	om.ForEachIndexed(func(idx int, k string, v string) {
		t.Logf("idx: %d, key : %s, value : %s", idx, k, v)
	})

	om.Delete("a")

	om.ForEachIndexed(func(idx int, k string, v string) {
		t.Logf("idx: %d, key : %s, value : %s", idx, k, v)
	})

	str := om.JoinToString(func(k string, v string) string {
		return fmt.Sprintf("{key:%s,v:%s}", k, v)
	})
	t.Logf(str)

}
