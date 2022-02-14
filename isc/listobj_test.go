package isc

import (
	"testing"
)

type MyStruct struct {
	Name string
	Age  int
}

func TestISCList_associateBy(t *testing.T) {
	var testList ISCList[MyStruct]
	s1 := MyStruct{
		Name: "K",
		Age:  1,
	}
	testList.Add(s1)
	testList.Add(MyStruct{
		Name: "K2",
		Age:  2,
	})
	testList.Add(MyStruct{
		Name: "K3",
		Age:  3,
	})
	l := AssociateBy[MyStruct](testList, func(t MyStruct) interface{} {
		return t.Name
	})
	println(ToString(l))
}
