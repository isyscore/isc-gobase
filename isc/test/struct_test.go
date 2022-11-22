package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/isc"
	"reflect"
	"testing"
)

func TestField(t *testing.T) {
	privateField := PrivateFieldStruct{}
	data := "data"
	isc.SetFieldPrivateValue(reflect.ValueOf(privateField), "name", reflect.ValueOf(&data))

	dataRel := isc.GetPrivateFieldValue(reflect.ValueOf(&privateField), "name")
	fmt.Println(dataRel)
}

type PrivateFieldStruct struct {
	name string
	age int
}
