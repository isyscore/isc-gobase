package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	. "github.com/isyscore/isc-gobase/isc"
)

type ApiPath struct {
	Path     string
	Handler  gin.HandlerFunc
	Method   HttpMethod
	Versions ISCList[*ApiVersion]
}

type ApiVersion struct {
	Header  string
	Version string
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
			if a.Version == c.GetHeader(a.Header) {
				a.Handler(c)
			}
		})
	}
	apiVersionList = append(apiVersionList, &v)
	return &v
}

func (ap *ApiPath) AddVersion(header string, version string, handler gin.HandlerFunc) {
	av := ap.Versions.Find(func(a *ApiVersion) bool {
		return a.Header == header && a.Version == version
	})
	if av == nil {
		a := &ApiVersion{Header: header, Version: version, Handler: handler}
		ap.Versions = append(ap.Versions, a)
	} else {
		panic(fmt.Sprintf("版本 %s-%s 已经存在", header, version))
	}
}
