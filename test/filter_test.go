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

type FS struct {
	Name  string
	Count int
}

func TestStructFilter(t *testing.T) {
	list := isc.NewListWithItems(
		FS{Name: "a", Count: 1},
		FS{Name: "b", Count: 2},
		FS{Name: "c", Count: 3},
		FS{Name: "d", Count: 4},
		FS{Name: "e", Count: 5},
	)

	l2 := list.Filter(func(item FS) bool {
		return item.Count%2 == 0
	})
	l2.ForEach(func(item FS) {
		t.Logf("%v\n", item)
	})

	b1 := list.Contains(FS{Name: "a", Count: 1})
	b2 := list.Contains(FS{Name: "x", Count: 66})
	t.Logf("%v %v\n", b1, b2)
}
