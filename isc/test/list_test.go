package test

import (
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestList(t *testing.T) {
	list := isc.NewList[string]()
	list.Add("aaaa")
	list.AddAll("bbbb", "ccc")
	t.Logf("list: %v\n", list)

	list.Clear()
	t.Logf("list: %v\n", list)

}

func TestMap(t *testing.T) {
	m := isc.NewMap[string, string]()
	m["aa"] = "bb"
	m.Put("cc", "dd")
	t.Logf("m: %v\n", m)
	m.Delete("cc")
	t.Logf("m: %v\n", m)
	m.Clear()
	t.Logf("m: %v\n", m)
}
