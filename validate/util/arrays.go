package util

import (
	"encoding/json"
)

func ArraysToString(dataArray []string) string {
	if len(dataArray) == 1 {
		return dataArray[0]
	}
	myValue, _ := json.Marshal(dataArray)
	return string(myValue)
}
