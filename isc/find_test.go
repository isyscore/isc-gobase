package isc

import (
	"fmt"
	"testing"
)

func TestIndexOf(t *testing.T) {
	list := []int{2, 4, 6, 9, 12}
	item := 10
	res := IndexOf[int](list, item)
	fmt.Println(res)
}
