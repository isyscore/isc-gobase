package test

import (
	"github.com/isyscore/isc-gobase/isc"
	"testing"
)

func TestString(t *testing.T) {
	var s isc.ISCString = "abcdefg"
	ss := s.Insert(3, "xyz")
	// ss := s.SubStringAfterLast(",")
	t.Logf("%v\n", ss) // abcxyzdefg

	sss := ss.Delete(3, 3)
	t.Logf("%v\n", sss) // abcdefg

}
