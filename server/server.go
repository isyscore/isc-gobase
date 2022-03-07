package server

import (
	"fmt"
	"io/ioutil"
	"strings"
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

func init() {

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
	logger.InitLog(level, timeFieldFormat, colored, appName)
}

func Run() {
	StartServer()
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

func RegisterWebSocketRoute(path string, svr *websocket.Server) {
	if !checkEngine() {
		return
	}
	engine.GET(path, svr.Handler())
}

func Post(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmPost, handler)
}

func PostApiModel(path string, handler gin.HandlerFunc) {
	Post(getPathAppendApiModel(path), handler)
}

func Delete(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmDelete, handler)
}

func DeleteApiModel(path string, handler gin.HandlerFunc) {
	Delete(getPathAppendApiModel(path), handler)
}

func Put(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmPut, handler)
}

func PutApiModel(path string, handler gin.HandlerFunc) {
	Put(getPathAppendApiModel(path), handler)
}

func Head(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmHead, handler)
}

func HeadApiModel(path string, handler gin.HandlerFunc) {
	Head(getPathAppendApiModel(path), handler)
}

func Get(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmGet, handler)
}

func GetApiModel(path string, handler gin.HandlerFunc) {
	Get(getPathAppendApiModel(path), handler)
}

func Options(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmOptions, handler)
}

func OptionsApiModel(path string, handler gin.HandlerFunc) {
	Options(getPathAppendApiModel(path), handler)
}

func GetPost(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmGetPost, handler)
}

func GetPostApiModel(path string, handler gin.HandlerFunc) {
	GetPost(getPathAppendApiModel(path), handler)
}

func All(path string, handler gin.HandlerFunc) {
	RegisterRoute(path, HmAll, handler)
}

func AllApiModel(path string, handler gin.HandlerFunc) {
	All(getPathAppendApiModel(path), handler)
}

func getPathAppendApiModel(path string) string {
	apiModel := config.GetValueString("api-module")
	if !strings.HasPrefix(apiModel, "/") {
		apiModel = "/" + apiModel
	}

	if !strings.HasSuffix(apiModel, "/") {
		apiModel = apiModel + "/"
	}

	if !strings.HasPrefix(path, "/") {
		return apiModel + path
	} else {
		return apiModel + path[:len(path)-1]
	}
}
