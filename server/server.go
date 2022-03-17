package server

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"

	"github.com/isyscore/isc-gobase/logger"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/websocket"
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

var ApiPrefix = "/api"

var engine *gin.Engine = nil

func InitServer() {

	isc.PrintBanner()

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

	engine = gin.New()
	engine.Use(Cors(), gin.Recovery())

	ap := config.GetValueStringDefault("base.api.prefix", "")
	if ap != "" {
		ApiPrefix = ap
	}

	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("base.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(ApiPrefix + "/" + config.GetValueString("api-module"))
	}

	// 注册 配置检测endpoint
	if config.GetValueBoolDefault("base.endpoint.config.enable", false) {
		RegisterConfigWatchEndpoint(ApiPrefix + "/" + config.GetValueString("api-module"))
	}
	level := config.GetValueStringDefault("server.logger.level", "info")
	timeFieldFormat := config.GetValueStringDefault("server.logger.time.format", time.RFC3339)
	colored := config.GetValueBoolDefault("server.logger.color.enable", false)
	appName := config.GetValueStringDefault("base.application.name", "isc-gobase")
	splitEnable := config.GetValueBoolDefault("server.logger.split.enable", false)
	splitSize := config.GetValueInt64Default("server.logger.split.size", 300)
	logger.InitLog(level, timeFieldFormat, colored, appName, splitEnable, splitSize)
}

func StartServer() {
	if !checkEngine() {
		return
	}
	logger.Info("开始启动服务")
	port := config.GetValueIntDefault("server.port", 8080)
	logger.Info("服务端口号: %d", port)
	err := engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error("启动服务异常 (%v)", err)
	}
}

func RegisterStatic(relativePath string, rootPath string) {
	if !checkEngine() {
		return
	}
	engine.Static(relativePath, rootPath)
}

func RegisterStaticFile(relativePath string, filePath string) {
	if !checkEngine() {
		return
	}
	engine.StaticFile(relativePath, filePath)
}

func RegisterPlugin(plugin gin.HandlerFunc) {
	if !checkEngine() {
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
		c.Data(200, "application/json; charset=utf-8", []byte(status()))
	})
	RegisterRoute(apiBase+"/system/init", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(init()))
	})
	RegisterRoute(apiBase+"/system/destroy", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(destroy()))
	})
}
func checkEngine() bool {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return false
	}
	return true
}
func RegisterRoute(path string, method HttpMethod, handler gin.HandlerFunc) {
	if !checkEngine() {
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

func RegisterRouteWithHeader(path string, method HttpMethod, header []string, versionName []string, handler gin.HandlerFunc) {
	if !checkEngine() {
		return
	}
	p := GetApiPath(path, method)
	if p == nil {
		p = NewApiPath(path, method)
		switch method {
		case HmAll:
			engine.GET(path, p.Handler)
			engine.POST(path, p.Handler)
			engine.PUT(path, p.Handler)
			engine.DELETE(path, p.Handler)
			engine.OPTIONS(path, p.Handler)
			engine.HEAD(path, p.Handler)
		case HmGet:
			engine.GET(path, p.Handler)
		case HmPost:
			engine.POST(path, p.Handler)
		case HmPut:
			engine.PUT(path, p.Handler)
		case HmDelete:
			engine.DELETE(path, p.Handler)
		case HmOptions:
			engine.OPTIONS(path, p.Handler)
		case HmHead:
			engine.HEAD(path, p.Handler)
		case HmGetPost:
			engine.GET(path, p.Handler)
			engine.POST(path, p.Handler)
		}
	}
	p.AddVersion(header, versionName, handler)
}

func RegisterWebSocketRoute(path string, svr *websocket.Server) {
	if !checkEngine() {
		return
	}
	engine.GET(path, svr.Handler())
}
