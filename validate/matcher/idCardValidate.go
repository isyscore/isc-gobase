package matcher

import (
	"regexp"
	"strconv"
)

var idCardSize = 17

// 加权因子
var weightFactor = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 校验码
var checkCode = [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

// 第一位不可能是0
// 第二位到第六位可以是0-9
// 第七位到第十位是年份，所以七八位为19或者20
// 十一位和十二位是月份，这两位是01-12之间的数值
// 十三位和十四位是日期，是从01-31之间的数值
// 十五，十六，十七都是数字0-9
// 十八位可能是数字0-9，也可能是X
var idCardPatter = "^[1-9][0-9]{5}([1][9][0-9]{2}|[2][0][0|1][0-9])([0][1-9]|[1][0|1|2])([0][1-9]|[1|2][0-9]|[3][0|1])[0-9]{3}([0-9]|[X])$"

func idCardIsValidate(idCard string) bool {
	if idCard == "" {
		return false
	}

	if len(idCard) < idCardSize {
		return false
	}

	seventeen := idCard[:idCardSize]

	num := 0
	for index, data := range seventeen {
		r, _ := strconv.Atoi(string(data))
		num += r * weightFactor[index]
	}

	if idCard[len(idCard)-1:] != checkCode[num%11] {
		return false
	}

	result, _ := regexp.MatchString(idCardPatter, idCard)
	return result
}
