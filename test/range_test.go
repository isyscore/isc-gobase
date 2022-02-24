package test

import (
	"github.com/isyscore/isc-gobase/ir"
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestRange(t *testing.T) {
	om := isc.NewOrderMap[string, string]()
	om.Put("a", "1")
	om.Put("b", "2")
	om.Put("c", "3")
	for _, item := range ir.OrderMap(om) {
		t.Logf("%s: %s", item.Key, item.Value)
	}

	for _, idx := range ir.IntStep(0, 10, 2) {
		t.Logf("%d", idx)
	}

}
