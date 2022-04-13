package test

import (
	"strconv"
	"testing"

	. "github.com/isyscore/isc-gobase/isc"
)

func TestNumber(t *testing.T) {
	i := ISCInt(1)
	iarr := i.RangeTo(10)
	iarr.ForEach(func(i int) {
		t.Logf("%d\n", i)
	})

	iarr2 := i.RangeStepTo(10, 2)
	iarr2.ForEach(func(i int) {
		t.Logf("%d\n", i)
	})

	ii := ISCInt(10)
	iiarr := ii.DownTo(0)
	iiarr.ForEach(func(i int) {
		t.Logf("%d\n", i)
	})
}

func TestRotate(t *testing.T) {
	a, b := ISCInt(1), ISCInt(8)
	s := strconv.IntSize
	t.Logf("int size = %d", strconv.IntSize)

	ii := a << 1
	iii := a >> (s - 1)

	// iii := a >> -1
	t.Logf("%d\n", ii)
	t.Logf("%d\n", iii)
	t.Logf("%d\n", ii|iii)
	// return ISCInt(uint((int(i) << bitCount)) | (uint(i) >> -bitCount))

	t.Logf("%d\n", a.RotateLeft(1))
	t.Logf("%d\n", a.RotateLeft(2))
	t.Logf("%d\n", b.RotateRight(1))
	t.Logf("%d\n", b.RotateRight(2))
}

func TestRadix(t *testing.T) {
	a := ISCInt(109)
	t.Logf("%s\n", a.ToHex())
	t.Logf("%s\n", a.ToOct())
	t.Logf("%s\n", a.ToBinary())

	b := ISCString("6D")
	i, _ := b.ToIntRadix(16)
	t.Logf("%d\n", i)

}

func TestRune(t *testing.T) {
	c := ISCChar('青')
	t.Logf("%d\n", c.Code())

	c2 := ISCChar(38738)
	t.Logf("%s\n", c2.ToString())
}

func TestForRange(t *testing.T) {
	for idx, e := range ISCInt(0).RangeStepTo(10, 2) {
		t.Logf("%d %d\n", idx, e)
	}
}
