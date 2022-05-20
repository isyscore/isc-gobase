package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type ValueModelIdCardEntity struct {
	Data string `match:"model=id_card"`
}

type ValueModelPhone struct {
	Data string `match:"model=phone"`
}

type ValueModelFixedPhoneEntity struct {
	Data string `match:"model=fixed_phone"`
}

type ValueModelEmailEntity struct {
	Data string `match:"model=mail"`
}

type ValueModelIpAddressEntity struct {
	Data string `match:"model=ip"`
}

// 身份证号
func TestModelIdCard(t *testing.T) {
	var value ValueModelIdCardEntity
	var result bool
	var err string

	// 测试 异常情况
	value = ValueModelIdCardEntity{Data: "4109281002226311"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 4109281002226311 不符合身份证要求", result, false)

	// 测试 异常情况
	value = ValueModelIdCardEntity{Data: "28712381"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 28712381 不符合身份证要求", result, false)
}

// 手机号
func TestModelPhone(t *testing.T) {
	var value ValueModelPhone
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelPhone{Data: "15700092345"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelPhone{Data: "28712381"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 28712381 没有命中只允许类型 [phone]", result, false)
}

// 固定电话
func TestModelFixedPhone(t *testing.T) {
	var value ValueModelFixedPhoneEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelFixedPhoneEntity{Data: "0393-3879765"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelFixedPhoneEntity{Data: "1387772"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 1387772 没有命中只允许类型 [fixed_phone]", result, false)
}

// 邮箱
func TestModelMail(t *testing.T) {
	var value ValueModelEmailEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelEmailEntity{Data: "123lan@163.com"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelEmailEntity{Data: "123@"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 123@ 没有命中只允许类型 [mail]", result, false)
}

// ip地址
func TestModelIpAddress(t *testing.T) {
	var value ValueModelIpAddressEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelIpAddressEntity{Data: "192.123.231.222"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelIpAddressEntity{Data: "123.231.222"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 123.231.222 没有命中只允许类型 [ip]", result, false)

	// 测试 异常情况
	value = ValueModelIpAddressEntity{Data: "192.123.231.adf"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 192.123.231.adf 没有命中只允许类型 [ip]", result, false)
}
