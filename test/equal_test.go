package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

func TestEqual(t *testing.T) {
	l1 := []int{1, 2, 3}
	l2 := []int{1, 2, 3}
	l3 := []int{1, 2, 4}
	b1 := isc.ListEquals(l1, l2)
	t.Logf("l1 == l2: %v", b1)
	b2 := isc.ListEquals(l1, l3)
	t.Logf("l1 == l3: %v", b2)

	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}
	b1 = isc.MapEquals(m1, m2)
	t.Logf("m1 == m2: %v", b1)
	b2 = isc.MapEquals(m1, m3)
	t.Logf("m1 == m3: %v", b2)
}
