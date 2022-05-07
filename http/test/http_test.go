package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/http"
	"testing"
)

func TestGetSimple(t *testing.T) {
	// {"code":"success","data":"ok","message":"成功"}
	data, err := http.GetSimple("http://localhost:8082/api/api/app/sample/test/get")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fmt.Println(string(data))

	// "ok"
	data, err = http.GetSimpleOfStandard("http://localhost:8082/api/api/app/sample/test/get")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fmt.Println(string(data))
}
