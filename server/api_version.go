package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	. "github.com/isyscore/isc-gobase/isc"
)

type ApiPath struct {
	Path     string
	Handler  gin.HandlerFunc
	Method   HttpMethod
	Versions ISCList[*ApiVersion]
}

type ApiVersion struct {
	Header  []string
	Version []string
	Handler gin.HandlerFunc
}

var apiVersionList ISCList[*ApiPath]

func GetApiPath(path string, method HttpMethod) *ApiPath {
	v := apiVersionList.Find(func(ap *ApiPath) bool {
		return ap.Path == path
	})
	if v == nil {
		return nil
	} else {
		return *v
	}
}

func NewApiPath(path string, method HttpMethod) *ApiPath {
	v := ApiPath{
		Path:   path,
		Method: method,
	}
	v.Handler = func(c *gin.Context) {
		v.Versions.ForEach(func(a *ApiVersion) {
			// 取出所有已定义header所对应的值
			t := isc.ListToMapFrom[string, string](a.Header).Map(func(item string) string {
				return c.GetHeader(item)
			})
			if isc.ListEquals(a.Version, t) {
				// 找到符合条件的路由版本，并转发请求
				a.Handler(c)
			}
		})
	}
	// 将路由添加到维护列表中，只有第一次添加时，会注册到gin
	apiVersionList = append(apiVersionList, &v)
	return &v
}

func (ap *ApiPath) AddVersion(header []string, version []string, handler gin.HandlerFunc) {
	// 查找指定版本的路由是否已经存在
	av := ap.Versions.Find(func(a *ApiVersion) bool {
		return isc.ListEquals(a.Header, header) && isc.ListEquals(a.Version, version)
	})
	if av == nil {
		// 不存在，则添加一个
		a := &ApiVersion{Header: header, Version: version, Handler: handler}
		ap.Versions = append(ap.Versions, a)
	} else {
		// 不允许重复添加
		panic(fmt.Sprintf("版本 %s-%s 已经存在", header, version))
	}
}
