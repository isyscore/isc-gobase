package test

import (
	"github.com/isyscore/isc-gobase/goid"
	"testing"
)

func TestUUID(t *testing.T) {
	id := goid.GenerateUUID()
	t.Logf("UUID: %s\n", id)
}
