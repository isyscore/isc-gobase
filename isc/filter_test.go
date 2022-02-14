package isc

import (
	"testing"
)

func TestListDistinct(t *testing.T) {
	list := []string{"1", "2", "test", "test", "7"}
	l := ListDistinct[string](list)

	println(ToString(l))
}
