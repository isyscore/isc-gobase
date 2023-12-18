package http

import (
	"context"
	"fmt"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/isc"
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

func Test_urlWithParameter(t *testing.T) {
	type args struct {
		url          string
		parameterMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal.",
			args: args{
				url: "http://127.0.0.1:38080",
				parameterMap: map[string]string{
					"msg": "{\"a\":1}",
				},
			},
			want: "http://127.0.0.1:38080?msg=%7B%22a%22%3A1%7D",
		},
		{
			name: "url with queryParam",
			args: args{
				url: "http://127.0.0.1:38080?example=1",
				parameterMap: map[string]string{
					"msg": "{\"a\":1}",
				},
			},
			want: "http://127.0.0.1:38080?example=1&msg=%7B%22a%22%3A1%7D",
		},
		{
			name: "url with simple",
			args: args{
				url: "www.example.com",
				parameterMap: map[string]string{
					"msg": "{\"a\":1}",
				},
			},
			want: "www.example.com?msg=%7B%22a%22%3A1%7D",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := baseHttp.UrlWithParameter(tt.args.url, tt.args.parameterMap); got != tt.want {
				t.Errorf("UrlWithParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
