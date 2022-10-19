package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/http"
	"testing"
)

func TestGetSimple(t *testing.T) {
	// {"code":"success","data":"ok","message":"成功"}
	_, _, data, err := http.GetSimple("http://localhost:8082/api/api/app/sample/test/get")
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	}
	fmt.Println(string(data.([]byte)))

	// "ok"
	_, _, data, err = http.GetSimpleOfStandard("http://localhost:8082/api/api/app/sample/test/get")
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	}
	fmt.Println(string(data.([]byte)))
}
