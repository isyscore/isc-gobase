package test

import (
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

type AssociateStruct struct {
	Key  string
	Name string
	Age  int
}

func initList() []AssociateStruct {
	return []AssociateStruct{
		{"K", "库陈胜", 20},
		{"K", "库陈胜", 30},
		{"K1", "库陈胜", 30},
		{"K1", "库陈胜", 40},
		{"K1", "库陈胜", 50},
		{"K2", "库陈胜", 60},
		{"K2", "库陈胜", 70},
		{"K2", "库陈胜", 80},
		{"K4", "库陈胜", 90},
	}
}

var transformFun = func(a AssociateStruct) isc.Pair[string, AssociateStruct] {
	return isc.Pair[string, AssociateStruct]{
		a.Key,
		a,
	}
}

var transformFun1 = func(a AssociateStruct) int {
	return a.Age
}

var keySelector = func(a AssociateStruct) string {
	return a.Key
}

func TestAssociate(t *testing.T) {
	list := initList()
	l := isc.Associate(list, transformFun)
	t.Logf("%v", l)
}

func TestAssociateTo(t *testing.T) {
	list := initList()
	m := map[string]AssociateStruct{}
	r := isc.AssociateTo(list, &m, transformFun)
	t.Logf("%v", r)
}

func TestAssociateBy(t *testing.T) {
	list := initList()
	r := isc.AssociateBy(list, keySelector)
	t.Logf("%v", r)
}

func TestAssociateByAndValue(t *testing.T) {
	list := initList()
	r := isc.AssociateByAndValue(list, keySelector, transformFun)
	t.Logf("%v", r)
}

func TestAssociateByAndValueTo(t *testing.T) {
	list := initList()
	m := make(map[string]int)
	isc.AssociateByAndValueTo(list, &m, keySelector, transformFun1)
	t.Logf("%v", m)
}

//
func TestAssociateByTo(t *testing.T) {
	list := initList()
	m := make(map[string]AssociateStruct)
	isc.AssociateByTo(list, &m, keySelector)
	t.Logf("%v", m)
}

//

//
func TestAssociateWith(t *testing.T) {
	list := initList()
	m := isc.AssociateWith(list, transformFun1)
	t.Logf("%v", m)
}

//
func TestAssociateWithTo(t *testing.T) {
	list := initList()
	m := make(map[AssociateStruct]int)
	isc.AssociateWithTo[AssociateStruct, int](list, &m, transformFun1)
	t.Logf("%v", m)
}
