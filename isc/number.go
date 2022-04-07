package isc

import (
	"fmt"
	"strconv"
)

type ISCInt int
type ISCInt8 int8
type ISCInt16 int16
type ISCInt32 int32
type ISCInt64 int64
type ISCChar uint8
type ISCFloat float32
type ISCFloat64 float64

func (i ISCInt) RangeTo(to int) ISCList[int] {
	var ret ISCList[int]
	for ii := int(i); ii <= to; ii++ {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt) RangeStepTo(to int, step int) ISCList[int] {
	var ret ISCList[int]
	for ii := int(i); ii <= to; ii += step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt) DownTo(to int) ISCList[int] {
	var ret ISCList[int]
	for ii := int(i); ii >= to; ii-- {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt) DownStepTo(to int, step int) ISCList[int] {
	var ret ISCList[int]
	for ii := int(i); ii >= to; ii -= step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt8) RangeTo(to int8) ISCList[int8] {
	var ret ISCList[int8]
	for ii := int8(i); ii <= to; ii++ {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt8) RangeStepTo(to int8, step int8) ISCList[int8] {
	var ret ISCList[int8]
	for ii := int8(i); ii <= to; ii += step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt8) DownTo(to int8) ISCList[int8] {
	var ret ISCList[int8]
	for ii := int8(i); ii >= to; ii-- {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt8) DownStepTo(to int8, step int8) ISCList[int8] {
	var ret ISCList[int8]
	for ii := int8(i); ii >= to; ii -= step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt16) RangeTo(to int16) ISCList[int16] {
	var ret ISCList[int16]
	for ii := int16(i); ii <= to; ii++ {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt16) RangeStepTo(to int16, step int16) ISCList[int16] {
	var ret ISCList[int16]
	for ii := int16(i); ii <= to; ii += step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt16) DownTo(to int16) ISCList[int16] {
	var ret ISCList[int16]
	for ii := int16(i); ii >= to; ii-- {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt16) DownStepTo(to int16, step int16) ISCList[int16] {
	var ret ISCList[int16]
	for ii := int16(i); ii >= to; ii -= step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt32) RangeTo(to int32) ISCList[int32] {
	var ret ISCList[int32]
	for ii := int32(i); ii <= to; ii++ {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt32) RangeStepTo(to int32, step int32) ISCList[int32] {
	var ret ISCList[int32]
	for ii := int32(i); ii <= to; ii += step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt32) DownTo(to int32) ISCList[int32] {
	var ret ISCList[int32]
	for ii := int32(i); ii >= to; ii-- {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt32) DownStepTo(to int32, step int32) ISCList[int32] {
	var ret ISCList[int32]
	for ii := int32(i); ii >= to; ii -= step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt64) RangeTo(to int64) ISCList[int64] {
	var ret ISCList[int64]
	for ii := int64(i); ii <= to; ii++ {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt64) RangeStepTo(to int64, step int64) ISCList[int64] {
	var ret ISCList[int64]
	for ii := int64(i); ii <= to; ii += step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt64) DownTo(to int64) ISCList[int64] {
	var ret ISCList[int64]
	for ii := int64(i); ii >= to; ii-- {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt64) DownStepTo(to int64, step int64) ISCList[int64] {
	var ret ISCList[int64]
	for ii := int64(i); ii >= to; ii -= step {
		ret = append(ret, ii)
	}
	return ret
}

func (i ISCInt) ToString() ISCString {
	return ISCString(fmt.Sprintf("%d", i))
}

func (i ISCInt8) ToString() ISCString {
	return ISCString(fmt.Sprintf("%d", i))
}

func (i ISCInt16) ToString() ISCString {
	return ISCString(fmt.Sprintf("%d", i))
}

func (i ISCInt32) ToString() ISCString {
	return ISCString(fmt.Sprintf("%d", i))
}

func (i ISCInt64) ToString() ISCString {
	return ISCString(fmt.Sprintf("%d", i))
}

func (i ISCChar) ToString() ISCString {
	return ISCString(string(i))
}

func (i ISCFloat) ToString() ISCString {
	return ISCString(fmt.Sprintf("%f", i))
}

func (i ISCFloat64) ToString() ISCString {
	return ISCString(fmt.Sprintf("%f", i))
}

func (i ISCInt) RotateLeft(bitCount int) ISCInt {
	s := strconv.IntSize
	return ISCInt(int(i)<<bitCount | int(i)>>(s-bitCount))
}

func (i ISCInt) RotateRight(bitCount int) ISCInt {
	s := strconv.IntSize
	return ISCInt(int(i)>>bitCount | int(i)<<(s-bitCount))
}

func (i ISCInt8) RotateLeft(bitCount int8) ISCInt8 {
	return ISCInt8(int8(i)<<bitCount | int8(i)>>(8-bitCount))
}

func (i ISCInt8) RotateRight(bitCount int8) ISCInt8 {
	return ISCInt8(int8(i)>>bitCount | int8(i)<<(8-bitCount))
}

func (i ISCInt16) RotateLeft(bitCount int16) ISCInt16 {
	return ISCInt16(int16(i)<<bitCount | int16(i)>>(16-bitCount))
}

func (i ISCInt16) RotateRight(bitCount int16) ISCInt16 {
	return ISCInt16(int16(i)>>bitCount | int16(i)<<(16-bitCount))
}

func (i ISCInt32) RotateLeft(bitCount int32) ISCInt32 {
	return ISCInt32(int32(i)<<bitCount | int32(i)>>(32-bitCount))
}

func (i ISCInt32) RotateRight(bitCount int32) ISCInt32 {
	return ISCInt32(int32(i)>>bitCount | int32(i)<<(32-bitCount))
}

func (i ISCInt64) RotateLeft(bitCount int64) ISCInt64 {
	return ISCInt64(int64(i)<<bitCount | int64(i)>>(64-bitCount))
}

func (i ISCInt64) RotateRight(bitCount int64) ISCInt64 {
	return ISCInt64(int64(i)>>bitCount | int64(i)<<(64-bitCount))
}

func (i ISCInt) ToHex() string {
	return fmt.Sprintf("%X", i)
}

func (i ISCInt8) ToHex() string {
	return fmt.Sprintf("%X", i)
}

func (i ISCInt16) ToHex() string {
	return fmt.Sprintf("%X", i)
}

func (i ISCInt32) ToHex() string {
	return fmt.Sprintf("%X", i)
}

func (i ISCInt64) ToHex() string {
	return fmt.Sprintf("%X", i)
}

func (i ISCInt) ToOct() string {
	return strconv.FormatInt(int64(i), 8)
}

func (i ISCInt8) ToOct() string {
	return strconv.FormatInt(int64(i), 8)
}

func (i ISCInt16) ToOct() string {
	return strconv.FormatInt(int64(i), 8)
}

func (i ISCInt32) ToOct() string {
	return strconv.FormatInt(int64(i), 8)
}

func (i ISCInt64) ToOct() string {
	return strconv.FormatInt(int64(i), 8)
}

func (i ISCInt) ToBinary() string {
	return strconv.FormatInt(int64(i), 2)
}

func (i ISCInt8) ToBinary() string {
	return strconv.FormatInt(int64(i), 2)
}

func (i ISCInt16) ToBinary() string {
	return strconv.FormatInt(int64(i), 2)
}

func (i ISCInt32) ToBinary() string {
	return strconv.FormatInt(int64(i), 2)
}

func (i ISCInt64) ToBinary() string {
	return strconv.FormatInt(int64(i), 2)
}

/*

===Char===
isLetter
isLetterOrDigit
isDigit
isIdentifierIgnorable
isISOControl
isWhitespace
isUpperCase
isLowerCase
uppercase
lowercase
isTitleCase
titlecase
isHighSurrogate
isLowSurrogate



*/
