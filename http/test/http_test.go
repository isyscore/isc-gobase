package test

import (
	"context"
	"fmt"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"net/http"
	"testing"
	"unsafe"
)

type DemoHttpHook struct {
}

func (*DemoHttpHook) Before(ctx context.Context, req *http.Request) context.Context {
	return ctx
}

func (*DemoHttpHook) After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error) {

}

func TestGetSimple(t *testing.T) {
	_, _, data, err := baseHttp.GetSimple("http://10.30.30.78:29013/api/core/license/osinfo")

	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	}
	fmt.Println("结果： " + string(data.([]byte)))

	datas := isc.ToInt(unsafe.Sizeof(data))

	fmt.Println("====" + isc.ToString(datas))
}
