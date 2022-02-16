package isc

import (
	"testing"
)

type GroupStruct struct {
	Key  string
	Name string
	Age  int
}

func TestGroupBy(t *testing.T) {
	list := []GroupStruct{
		{Key: "K", Name: "库陈胜", Age: 1},
		{Key: "K", Name: "库陈胜1", Age: 2},
		{Key: "K", Name: "库陈胜2", Age: 3},
		{Key: "K1", Name: "库陈胜", Age: 1},
		{Key: "K1", Name: "库陈胜2", Age: 2},
		{Key: "K2", Name: "库陈胜3", Age: 1},
	}
	m := GroupBy[GroupStruct](list, func(t GroupStruct) string {
		return t.Key
	})
	t.Logf("%v", m)
}

func TestGroupByTransform(t *testing.T) {
	list := []GroupStruct{
		{Key: "K", Name: "库陈胜", Age: 1},
		{Key: "K", Name: "库陈胜1", Age: 2},
		{Key: "K", Name: "库陈胜2", Age: 3},
		{Key: "K1", Name: "库陈胜", Age: 1},
		{Key: "K1", Name: "库陈胜2", Age: 2},
		{Key: "K2", Name: "库陈胜3", Age: 1},
	}
	m := GroupByTransform[GroupStruct](list, func(t GroupStruct) string {
		return t.Key
	}, func(t GroupStruct) int {
		return t.Age
	})
	t.Logf("%v", m)
}
