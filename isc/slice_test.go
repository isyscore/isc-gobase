package isc

import (
	"testing"
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
	l := SliceDistinctTo[sliceTestStruct](list, func(s sliceTestStruct) string {
		return s.Name
	})
	println(ToString(l))
	b := Contains[sliceTestStruct](list, func(s sliceTestStruct) string {
		return s.Name
	}, "库陈胜")
	println(b)
}
