package fun

import (
	"fmt"
	"github.com/isyscore/isc-gobase/validate"
)

type CustomizeEntity1 struct {
	Name string `match:"customize=judge1Name"`
}

type CustomizeEntity2 struct {
	Name string `match:"customize=judge2Name"`
}

type CustomizeEntity3 struct {
	Name string `match:"customize=judge3Name"`
	Age  int
}

type CustomizeEntity4 struct {
	Name string `match:"customize=judge4Name"`
	Age  int
}

type CustomizeEntity5 struct {
	Name string `match:"customize=judge5Name"`
	Age  int
}

type CustomizeEntity6 struct {
	Name *string `match:"customize=judge6Name"`
	Flag *bool `match:"customize=judge6Flag"`
	Flag2 *bool `match:"customize=judge6Flag2"`
	Age  int
}

func JudgeString1(name string) bool {
	if name == "zhou" || name == "宋江" {
		return true
	}

	return false
}

func JudgeString2(name string) (bool, string) {
	if name == "zhou" || name == "宋江" {
		return true, ""
	}

	return false, "没有命中可用的值'zhou'和'宋江'"
}

func JudgeString3(customize CustomizeEntity3, name string) (bool, string) {
	if name == "zhou" || name == "宋江" {
		if customize.Age > 12 {
			return true, ""
		} else {
			return false, "用户[" + name + "]" + "没有满足年龄age > 12，" + "当前年龄为：" + fmt.Sprintf("%v", customize.Age)
		}

	} else {
		return false, "没有命中可用的值'zhou'和'宋江'"
	}
}
func JudgeString4(customize CustomizeEntity4, name string) (string, bool) {
	if name == "zhou" || name == "宋江" {
		if customize.Age > 12 {
			return "", true
		} else {
			return "用户[" + name + "]" + "没有满足年龄age > 12，" + "当前年龄为：" + fmt.Sprintf("%v", customize.Age), false
		}

	} else {
		return "没有命中可用的值'zhou'和'宋江'", false
	}
}

func JudgeString5(customize CustomizeEntity5) (string, bool) {
	var name = customize.Name
	if name == "zhou" || name == "宋江" {
		if customize.Age > 12 {
			return "", true
		} else {
			return "用户[" + name + "]" + "没有满足年龄age > 12，" + "当前年龄为：" + fmt.Sprintf("%v", customize.Age), false
		}

	} else {
		return "没有命中可用的值'zhou'和'宋江'", false
	}
}

func JudgeString6(customize CustomizeEntity6) (string, bool) {
	var name = customize.Name
	if name == nil {
		return "用户不可为空", false
	}
	return "", true
}

func JudgeString6Flag(customize CustomizeEntity6) (string, bool) {
	var flag = customize.Flag
	if flag == nil {
		return "flag不可为空", false
	}
	return "", true
}

func JudgeString6Flag2(flag *bool) (string, bool) {
	if flag == nil {
		return "flag不可为空", false
	}
	return "", true
}

func init() {
	validate.RegisterCustomize("judge1Name", JudgeString1)
	validate.RegisterCustomize("judge2Name", JudgeString2)
	validate.RegisterCustomize("judge3Name", JudgeString3)
	validate.RegisterCustomize("judge4Name", JudgeString4)
	validate.RegisterCustomize("judge5Name", JudgeString5)
	validate.RegisterCustomize("judge6Name", JudgeString6)
	validate.RegisterCustomize("judge6Flag", JudgeString6Flag)
	validate.RegisterCustomize("judge6Flag2", JudgeString6Flag2)
}
