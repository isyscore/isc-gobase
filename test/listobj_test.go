package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

type MyStruct struct {
	Name string
	Age  int
}

func TestISCList_associateBy(t *testing.T) {
	var testList isc.ISCList[MyStruct]
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
	l := isc.AssociateBy(testList, func(t MyStruct) interface{} {
		return t.Name
	})
	t.Logf("%v\n", isc.ToString(l))
}
