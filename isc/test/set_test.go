package test

import (
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestSet(t *testing.T) {
	l := isc.NewListWithItems(1, 2, 3, 4, 1, 2, 3, 4, 5, 6)
	t.Logf("%v", l)
	ls := isc.ListToSet(l)
	t.Logf("%v", ls)

	s := isc.NewSetWithItems(1, 2, 3, 4, 1, 2, 3, 4, 5, 6)
	t.Logf("%v", s)

	_ = s.Add(7)
	s.AddAll(8, 9)
	t.Logf("%v", s)
	_ = s.Delete(5)
	t.Logf("%v", s)
	s.Clear()
	t.Logf("%v", s)
}
