package server

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"io/ioutil"
	"net/http"

	h2 "github.com/isyscore/isc-gobase/http"
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

var engine *gin.Engine = nil

func InitServer() {
	config.LoadConfig()

	mode := config.GetValueString("server.gin.mode")
	if "debug" == mode {
		gin.SetMode(gin.DebugMode)
	} else if "test" == mode {
		gin.SetMode(gin.TestMode)
	} else if "release" == mode {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	engine = gin.Default()
	engine.Use(Cors())

	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("base.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(config.GetValueString("api-module"))
	}

	// 注册 配置检测endpoint
	if config.GetValueBoolDefault("base.endpoint.config.enable", false) {
		RegisterConfigWatchEndpoint(config.GetValueString("api-module"))
	}
	logger.InitLog()
}

func StartServer() {
	if engine != nil {
		logger.Info("启动服务 ...")
		fmt.Println(config.GetValueString("server.port"))
		err := engine.Run(fmt.Sprintf(":%d", config.GetValueIntDefault("server.port", 8080)))
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

func RegisterHealthCheckEndpoint(apiBase string) {
	if "" == apiBase {
		return
	}
	RegisterRoute(apiBase+"/system/status", HmAll, healthSystemStatus)
	RegisterRoute(apiBase+"/system/init", HmAll, healthSystemInit)
	RegisterRoute(apiBase+"/system/destroy", HmAll, healthSystemDestroy)
}

func RegisterConfigWatchEndpoint(apiBase string) {
	if "" == apiBase {
		return
	}
	RegisterRoute(apiBase+"/config/values", HmGet, config.GetConfigValues)
	RegisterRoute(apiBase+"/config/value/:key", HmGet, config.GetConfigValue)
	RegisterRoute(apiBase+"/config/update", HmPut, config.UpdateConfig)
}

func RegisterCustomHealthCheck(apiBase string, status func() string, init func() string, destroy func() string) {
	RegisterRoute(apiBase+"/system/status", HmAll, func(c *gin.Context) {
		c.Data(http.StatusOK, h2.ContentTypeJson, []byte(status()))
	})
	RegisterRoute(apiBase+"/system/init", HmAll, func(c *gin.Context) {
		c.Data(http.StatusOK, h2.ContentTypeText, []byte(init()))
	})
	RegisterRoute(apiBase+"/system/destroy", HmAll, func(c *gin.Context) {
		c.Data(http.StatusOK, h2.ContentTypeText, []byte(destroy()))
	})
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
