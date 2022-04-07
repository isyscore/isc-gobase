package isc

import (
	"fmt"
	"strconv"
	"unicode"
)

type ISCInt int
type ISCInt8 int8
type ISCInt16 int16
type ISCInt32 int32
type ISCInt64 int64
type ISCFloat float32
type ISCFloat64 float64
type ISCChar rune

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

func (i ISCChar) RangeTo(to rune) ISCList[rune] {
	var ret ISCList[rune]
	for ii := int64(i); ii <= int64(to); ii++ {
		ret = append(ret, rune(ii))
	}
	return ret
}

func (i ISCChar) RangeStepTo(to rune, step int64) ISCList[rune] {
	var ret ISCList[rune]
	for ii := int64(i); ii <= int64(to); ii += step {
		ret = append(ret, rune(ii))
	}
	return ret
}

func (i ISCChar) DownTo(to rune) ISCList[rune] {
	var ret ISCList[rune]
	for ii := int64(i); ii >= int64(to); ii-- {
		ret = append(ret, rune(ii))
	}
	return ret
}

func (i ISCChar) DownStepTo(to rune, step int64) ISCList[rune] {
	var ret ISCList[rune]
	for ii := int64(i); ii >= int64(to); ii -= step {
		ret = append(ret, rune(ii))
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

func (i ISCChar) Code() int {
	return int(i)
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

func (i ISCChar) IsLetter() bool {
	return unicode.IsLetter(rune(i))
}

func (i ISCChar) IsDigit() bool {
	return unicode.IsDigit(rune(i))
}

func (i ISCChar) IsLetterOrDigit() bool {
	return unicode.IsLetter(rune(i)) || unicode.IsDigit(rune(i))
}

func (i ISCChar) IsSymbol() bool {
	return unicode.IsSymbol(rune(i))
}

func (i ISCChar) IsWhitespace() bool {
	return unicode.IsSpace(rune(i))
}

func (i ISCChar) ToUpper() ISCChar {
	return ISCChar(unicode.ToUpper(rune(i)))
}

func (i ISCChar) ToLower() ISCChar {
	return ISCChar(unicode.ToLower(rune(i)))
}

func (i ISCChar) ToTitle() ISCChar {
	return ISCChar(unicode.ToTitle(rune(i)))
}

func (i ISCChar) IsUpper() bool {
	return unicode.IsUpper(rune(i))
}

func (i ISCChar) IsLower() bool {
	return unicode.IsLower(rune(i))
}

func (i ISCChar) IsTitle() bool {
	return unicode.IsTitle(rune(i))
}

func (i ISCChar) IsISOControl() bool {
	return unicode.IsControl(rune(i))
}

const (
	MIN_LOW_SURROGATE  = 0xDC00
	MAX_LOW_SURROGATE  = 0xDFFF
	MIN_HIGH_SURROGATE = 0xD800
	MAX_HIGH_SURROGATE = 0xDBFF
)

func (i ISCChar) IsHighSurrogate() bool {
	return i >= MIN_HIGH_SURROGATE && i < (MAX_HIGH_SURROGATE+1)
}

func (i ISCChar) IsLowSurrogate() bool {
	return i >= MIN_LOW_SURROGATE && i < (MAX_LOW_SURROGATE+1)
}
