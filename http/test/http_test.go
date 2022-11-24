package test

import (
	"context"
	"fmt"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"net/http"
	"testing"
)

type DemoHttpHook struct {
}

func (*DemoHttpHook) Before(ctx context.Context, req *http.Request) context.Context {

	fmt.Println("信息" + req.URL.Path)
	return ctx
}

func (*DemoHttpHook) After(ctx context.Context, req *http.Request, res *http.Response, err error) {
	fmt.Println("返回值")
}

func init() {
	baseHttp.AddHook(&DemoHttpHook{})
}

func TestGetSimple(t *testing.T) {
	// {"code":"success","data":"ok","message":"成功"}
	baseHttp.GetSimple("http://10.30.30.78:29013/api/core/license/osinfo")

	//if err != nil {
	//	fmt.Printf("error = %v\n", err)
	//	return
	//}
	//fmt.Println(string(data.([]byte)))
	//
	//// "ok"
	//_, _, data, err = baseHttp.GetSimpleOfStandard("http://localhost:8082/api/api/app/sample/test/get")
	//if err != nil {
	//	fmt.Printf("error = %v\n", err)
	//	return
	//}
	//fmt.Println(string(data.([]byte)))
}
