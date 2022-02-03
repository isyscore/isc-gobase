package server

import (
	"fmt"

	"github.com/isyscore/isc-gobase/logger"

	"github.com/gin-gonic/gin"
)

type HttpMethod int

const (
	HmAll HttpMethod = iota
	HmGet
	HmPost
	HmPut
	HmDelete
	HmOptions
	HmHead
	HmGetPost
)

var serverPort = 8080
var engine *gin.Engine = nil

func InitServer(port int) {
	serverPort = port
	gin.SetMode(gin.ReleaseMode)
	engine = gin.Default()
	engine.Use(Cors())
}

func StartServer() {
	if engine != nil {
		logger.Info("启动服务 ...")
		err := engine.Run(fmt.Sprintf(":%d", serverPort))
		if err != nil {
			logger.Error("启动服务异常 (%v)", err)
		}
	} else {
		logger.Error("服务没有初始化，请先调用 InitServer")
	}
}

func RegisterStatic(relativePath string, rootPath string) {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return
	}
	engine.Static(relativePath, rootPath)
}

func RegisterStaticFile(relativePath string, filePath string) {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return
	}
	engine.StaticFile(relativePath, filePath)
}

func RegisterPlugin(plugin gin.HandlerFunc) {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return
	}
	engine.Use(plugin)
}

func Engine() *gin.Engine {
	return engine
}

func RegisterRoute(path string, method HttpMethod, handler gin.HandlerFunc) {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return
	}
	switch method {
	case HmAll:
		engine.GET(path, handler)
		engine.POST(path, handler)
		engine.PUT(path, handler)
		engine.DELETE(path, handler)
		engine.OPTIONS(path, handler)
		engine.HEAD(path, handler)
	case HmGet:
		engine.GET(path, handler)
	case HmPost:
		engine.POST(path, handler)
	case HmPut:
		engine.PUT(path, handler)
	case HmDelete:
		engine.DELETE(path, handler)
	case HmOptions:
		engine.OPTIONS(path, handler)
	case HmHead:
		engine.HEAD(path, handler)
	case HmGetPost:
		engine.GET(path, handler)
		engine.POST(path, handler)
	}

}
