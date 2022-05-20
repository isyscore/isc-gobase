package test

import (
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestGeneric(t *testing.T) {
	list := isc.NewListWithItems(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	i := list.IndexOf(5)
	t.Logf("indexOf(5) = %d\n", i)

	l2 := list.Filter(func(item int) bool {
		return item%2 == 0
	})

	t.Logf("filter = %v\n", l2)

	lg := isc.ListToTripleFrom[int, string, string](list)
	mg := lg.GroupBy(func(item int) string {
		if item%2 == 0 {
			return "even"
		} else {
			return "odd"
		}
	})
	t.Logf("grouped = %v\n", mg)
}
