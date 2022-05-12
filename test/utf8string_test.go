package test

import (
	"testing"

	. "github.com/isyscore/isc-gobase/isc"
)

func TestUTF8String(t *testing.T) {
	var s = ISCUTF8String("指令集UTF8字符串")

	t.Logf("len = %d\n", s.Length())
	idx := s.IndexOf(ISCUTF8String("集")) // 2
	t.Logf("indexOf(\"集\") = %d\n", idx)

	ss := s.Insert(3, ISCUTF8String("xyz"))
	t.Logf("%v\n", ss) // 指令集xyzUTF8字符串

	sss := ss.Delete(3, 3)
	t.Logf("%v\n", sss) // 指令集UTF8字符串

}
