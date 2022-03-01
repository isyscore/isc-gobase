package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

func TestListDistinct(t *testing.T) {
	list := []string{"1", "2", "test", "test", "7"}
	l := isc.ListDistinct(list)
	t.Logf("%v\n", l)
}

func TestFilter(t *testing.T) {
	list := isc.NewListWithItems(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	even := list.Filter(func(item int) bool {
		return item%2 == 0
	})
	even.ForEach(func(item int) {
		t.Logf("%v\n", item)
	})
}
