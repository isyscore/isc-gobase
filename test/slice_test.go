package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

type sliceTestStruct struct {
	Name string
	Age  int
}

func TestSliceDistinctTo(t *testing.T) {
	s1 := sliceTestStruct{
		Name: "库陈胜",
		Age:  30,
	}
	s2 := sliceTestStruct{
		Name: "酷达舒",
		Age:  29,
	}
	s3 := sliceTestStruct{
		Name: "库陈胜",
		Age:  28,
	}
	list := []sliceTestStruct{s1, s2, s3}
	l := isc.SliceDistinctTo(list, func(s sliceTestStruct) string {
		return s.Name
	})
	t.Logf("%s\n", isc.ToString(l))
	b := isc.SliceContains(list, func(s sliceTestStruct) string {
		return s.Name
	}, "库陈胜")
	t.Logf("%v\n", b)
}
