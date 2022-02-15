package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

func TestOperator(t *testing.T) {
	list := isc.NewListWithItems(1, 2, 3, 4, 5, 6, 7)
	l2 := isc.ListPlus(list, []int{7, 8, 9, 10})
	t.Logf("%v\n", l2)

	l3 := isc.ListMinus(list, []int{1, 3, 7, 8, 9})
	t.Logf("%v\n", l3)

	m1 := isc.NewMap[int, string]()
	m1[1] = "a"
	m1[2] = "b"
	m1[3] = "c"

	mt := isc.NewMap[int, string]()
	mt[1] = "a"
	mt[2] = "b"
	mt[4] = "d"
	mt[5] = "e"

	m2 := isc.MapPlus(m1, mt)
	t.Logf("%v\n", m2)
	m3 := isc.MapMinus(m1, mt)
	t.Logf("%v\n", m3)

}
