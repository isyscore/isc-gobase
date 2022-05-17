package test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/isyscore/isc-gobase/isc"
)

func TestStringConvert(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.BigCamel(originalStr)
	assert.Equal(t, "DataBaseUser", newStr)
}

func TestMiddleLine(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.MiddleLine(originalStr)
	assert.Equal(t, "data-base-user", newStr)
}

func TestBigCamelToMiddleLine(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToMiddleLine(originalStr)
	assert.Equal(t, "data-base-user", newStr)
}

func TestBigCamelToSmallCamel(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToSmallCamel(originalStr)
	assert.Equal(t, "dataBaseUser", newStr)
}

func TestBigCamelToPostUnder(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToPostUnder(originalStr)
	assert.Equal(t, "data_base_user_", newStr)
}

func TestPostUnder(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.PostUnder(originalStr)
	assert.Equal(t, "data_base_user_", newStr)
}

func TestPrePostUnder(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.PrePostUnder(originalStr)
	assert.Equal(t, "_data_base_user_", newStr)
}

func TestBigCamelToPrePostUnder(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToPrePostUnder(originalStr)
	assert.Equal(t, "_data_base_user_", newStr)
}

func TestPreUnder(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.PreUnder(originalStr)
	assert.Equal(t, "_data_base_user", newStr)
}

func TestBigCamelToPreUnder(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToPreUnder(originalStr)
	assert.Equal(t, "_data_base_user", newStr)
}

func TestBigCamelToUnderLine(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToUnderLine(originalStr)
	assert.Equal(t, "data_base_user", newStr)
}

func TestBigCamelToUpperMiddle(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToUpperMiddle(originalStr)
	assert.Equal(t, "DATA-BASE-USER", newStr)
}

func TestUpperUnderMiddle(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.UpperUnderMiddle(originalStr)
	assert.Equal(t, "DATA-BASE-USER", newStr)
}

func TestUpperUnder(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.UpperUnder(originalStr)
	assert.Equal(t, "DATA_BASE_USER", newStr)
}

func TestBigCamelToUpperUnder(t *testing.T) {
	originalStr := "DataBaseUser"
	newStr := isc.BigCamelToUpperUnder(originalStr)
	assert.Equal(t, "DATA_BASE_USER", newStr)
}

func TestMiddleLineToSmallCamel(t *testing.T) {
	originalStr := "data-base-user"
	newStr := isc.MiddleLineToSmallCamel(originalStr)
	assert.Equal(t, "dataBaseUser", newStr)
}

func TestMiddleLineToBigCamel(t *testing.T) {
	originalStr := "data-base-user"
	newStr := isc.MiddleLineToBigCamel(originalStr)
	assert.Equal(t, "DataBaseUser", newStr)
}

func TestPreFixUnderLine(t *testing.T) {
	originalStr := "dataBaseUser"
	newStr := isc.PreFixUnderLine(originalStr, "pre_")
	assert.Equal(t, "pre_data_base_user", newStr)
}

func TestUnderLineToSmallCamel(t *testing.T) {
	originalStr1 := "data_base_user"
	newStr1 := isc.UnderLineToSmallCamel(originalStr1)
	assert.Equal(t, "dataBaseUser", newStr1)

	originalStr2 := "_data_base_user"
	newStr2 := isc.UnderLineToSmallCamel(originalStr2)
	assert.Equal(t, "dataBaseUser", newStr2)

	originalStr3 := "data_base_user_"
	newStr3 := isc.UnderLineToSmallCamel(originalStr3)
	assert.Equal(t, "dataBaseUser", newStr3)
}

func TestPreFixUnderToSmallCamel(t *testing.T) {
	originalStr := "pre_data_base_user"
	newStr := isc.PreFixUnderToSmallCamel(originalStr, "pre_")
	assert.Equal(t, "dataBaseUser", newStr)
}

func TestUnderLineToBigCamel(t *testing.T) {
	originalStr1 := "data_base_user"
	newStr1 := isc.UnderLineToBigCamel(originalStr1)
	assert.Equal(t, "DataBaseUser", newStr1)

	originalStr2 := "_data_base_user"
	newStr2 := isc.UnderLineToBigCamel(originalStr2)
	assert.Equal(t, "DataBaseUser", newStr2)

	originalStr3 := "_data_base_user_"
	newStr3 := isc.UnderLineToBigCamel(originalStr3)
	assert.Equal(t, "DataBaseUser", newStr3)

	originalStr4 := "data_base_user_"
	newStr4 := isc.UnderLineToBigCamel(originalStr4)
	assert.Equal(t, "DataBaseUser", newStr4)
}

func TestUpperUnderMiddleToSmallCamel(t *testing.T) {
	originalStr := "DATA-BASE-USER"
	newStr := isc.UpperUnderMiddleToSmallCamel(originalStr)
	assert.Equal(t, "dataBaseUser", newStr)
}

func TestUpperUnderToSmallCamel(t *testing.T) {
	originalStr := "DATA_BASE_USER"
	newStr := isc.UpperUnderToSmallCamel(originalStr)
	assert.Equal(t, "dataBaseUser", newStr)
}

func TestUpperUnderToBigCamel(t *testing.T) {
	originalStr := "DATA_BASE_USER"
	newStr := isc.UpperUnderToBigCamel(originalStr)
	assert.Equal(t, "DataBaseUser", newStr)
}

func TestUpperMiddleToBigCamel(t *testing.T) {
	originalStr := "DATA-BASE-USER"
	newStr := isc.UpperMiddleToBigCamel(originalStr)
	assert.Equal(t, "DataBaseUser", newStr)
}
