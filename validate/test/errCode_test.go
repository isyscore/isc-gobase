package main

import (
	"github.com/isyscore/isc-gobase/validate"
	"testing"
)

type ErrCodeEntity1 struct {
	Name string `match:"value={zhou, song}" errCode:"12"`
}

type ErrCodeEntity2 struct {
	Name string `match:"customize=judgeErrCode1" errCode:"12"`
	Name2 string `match:"customize=judgeErrCode2" errCode:"12"`
}

func TestErrCode1(t *testing.T) {
	entity1 := ErrCodeEntity1{Name: "chen"}
	result, errCode, _ := validate.Check(entity1, "name")
	False(t, result)
	Equal(t, errCode, "12")

	entity1 = ErrCodeEntity1{Name: "zhou"}
	result, errCode, _ = validate.Check(entity1, "name")
	True(t, result)
	Equal(t, errCode, "")
}

func TestErrCode2(t *testing.T) {
	entity1 := ErrCodeEntity2{Name: "chen"}
	result, errCode, _ := validate.Check(entity1, "name")
	False(t, result)
	Equal(t, errCode, "123")

	entity2 := ErrCodeEntity2{Name2: "chen"}
	result, errCode, _ = validate.Check(entity2, "name2")
	False(t, result)
	Equal(t, errCode, "12")

	entity1 = ErrCodeEntity2{Name: "zhou"}
	result, errCode, _ = validate.Check(entity1, "name")
	True(t, result)
	Equal(t, errCode, "")
}

func judgeErrCode1(name string) (bool, string, string) {
	if name != "zhou" && name != "song" {
		return false, "123", "只支持song和zhou"
	}
	return true, "", ""
}

func judgeErrCode2(name string) (bool, string, string) {
	if name != "zhou" && name != "song" {
		return false, "", "只支持song和zhou"
	}
	return true, "", ""
}

func init()  {
	validate.RegisterCustomize("judgeErrCode1", judgeErrCode1)
	validate.RegisterCustomize("judgeErrCode2", judgeErrCode2)
}


