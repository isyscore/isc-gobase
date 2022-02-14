package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/isc"
)

func TestListDistinct(t *testing.T) {
	list := []string{"1", "2", "test", "test", "7"}
	l := isc.ListDistinct(list)
	t.Logf("%v\n", l)
}
