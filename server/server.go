package server

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/listener"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/isyscore/isc-gobase/server/rsp"

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

var GoBaseVersion = "1.1.0"
var ApiPrefix = "/api"

var engine *gin.Engine = nil

func init() {
	isc.PrintBanner()
	config.LoadConfig()
	printVersionAndProfile()

	if config.ExistConfigFile() && config.GetValueBoolDefault("base.server.enable", false) {
		InitServer()
	}
}

func InitServer() {
	if !config.ExistConfigFile() {
		logger.Error("没有找到任何配置文件，服务启动失败")
		return
	}
	mode := config.GetValueStringDefault("base.server.gin.mode", "release")
	if "debug" == mode {
		gin.SetMode(gin.DebugMode)
	} else if "test" == mode {
		gin.SetMode(gin.TestMode)
	} else if "release" == mode {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	engine = gin.New()
	engine.Use(Cors(), gin.Recovery())
	engine.Use(rsp.ResponseHandler())

	ap := config.GetValueStringDefault("base.api.prefix", "")
	if ap != "" {
		ApiPrefix = ap
	}

	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("base.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(ApiPrefix + "/" + config.ApiModule)
	}

	// 注册 配置检测endpoint
	if config.GetValueBoolDefault("base.endpoint.config.enable", false) {
		RegisterConfigWatchEndpoint(ApiPrefix + "/" + config.ApiModule)
	}

	// 注册 bean管理的功能
	if config.GetValueBoolDefault("base.endpoint.bean.enable", false) {
		RegisterBeanWatchEndpoint(ApiPrefix + "/" + config.ApiModule)
	}

	appName := config.GetValueStringDefault("base.application.name", "isc-gobase")

	var loggerCfg logger.LoggerConfig
	if err := config.GetValueObject("base.logger", &loggerCfg); err != nil {
		logger.Warn("获取配置失败", err)
	} else {
		logger.InitLog(appName, &loggerCfg)
	}
}

func printVersionAndProfile() {
	fmt.Printf("----------------------------- isc-gobase: %s --------------------------\n", GoBaseVersion)
	fmt.Printf("profile：%s\n", config.CurrentProfile)
	fmt.Printf("--------------------------------------------------------------------------\n")
}

func Run() {
	StartServer()
}

func StartServer() {
	if !checkEngine() {
		return
	}

	if !config.GetValueBoolDefault("base.server.enable", true) {
		return
	}

	logger.Info("开始启动服务")
	port := config.GetValueIntDefault("base.server.port", 8080)
	logger.Info("服务端口号: %d", port)

	graceRun(port)
}

func graceRun(port int) {
	engineServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: engine}
	go func() {
		if err := engineServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("启动服务异常 (%v)", err)
		} else {
			// 发送服务关闭事件
			listener.PublishEvent(listener.ServerStopEvent{})
		}
	}()

	// 发送服务启动事件
	listener.PublishEvent(listener.ServerFinishEvent{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("服务端准备关闭...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := engineServer.Shutdown(ctx); err != nil {
		logger.Warn("服务关闭异常: ", err)
	}
	logger.Info("服务端退出")
}

func RegisterStatic(relativePath string, rootPath string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.Static(relativePath, rootPath)
	return engine
}

func RegisterStaticFile(relativePath string, filePath string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.StaticFile(relativePath, filePath)
	return engine
}

func RegisterPlugin(plugin gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.Use(plugin)
	return engine
}

func Engine() gin.IRoutes {
	return engine
}

func RegisterHealthCheckEndpoint(apiBase string) gin.IRoutes {
	if "" == apiBase {
		return nil
	}
	RegisterRoute(apiBase+"/system/status", HmAll, healthSystemStatus)
	RegisterRoute(apiBase+"/system/init", HmAll, healthSystemInit)
	RegisterRoute(apiBase+"/system/destroy", HmAll, healthSystemDestroy)
	return engine
}

func RegisterConfigWatchEndpoint(apiBase string) gin.IRoutes {
	if "" == apiBase {
		return nil
	}
	RegisterRoute(apiBase+"/config/values", HmGet, config.GetConfigValues)
	RegisterRoute(apiBase+"/config/value/:key", HmGet, config.GetConfigValue)
	RegisterRoute(apiBase+"/config/update", HmPut, config.UpdateConfig)
	return engine
}

func RegisterBeanWatchEndpoint(apiBase string) gin.IRoutes {
	if "" == apiBase {
		return nil
	}
	RegisterRoute(apiBase+"/bean/name/all", HmGet, bean.DebugBeanAll)
	RegisterRoute(apiBase+"/bean/name/list/:name", HmGet, bean.DebugBeanList)
	RegisterRoute(apiBase+"/bean/field/get", HmPost, bean.DebugBeanGetField)
	RegisterRoute(apiBase+"/bean/field/set", HmPut, bean.DebugBeanSetField)
	RegisterRoute(apiBase+"/bean/fun/call", HmPost, bean.DebugBeanFunCall)
	return engine
}

func RegisterCustomHealthCheck(apiBase string, status func() string, init func() string, destroy func() string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	RegisterRoute(apiBase+"/system/status", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(status()))
	})
	RegisterRoute(apiBase+"/system/init", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(init()))
	})
	RegisterRoute(apiBase+"/system/destroy", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(destroy()))
	})
	return engine
}

func checkEngine() bool {
	if engine == nil {
		logger.Error("服务没有初始化，请先调用 InitServer")
		return false
	}
	return true
}

func RegisterRoute(path string, method HttpMethod, handler gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
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
	return engine
}

func RegisterRouteWithHeaders(path string, method HttpMethod, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
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
	return engine
}

func RegisterWebSocketRoute(path string, svr *websocket.Server) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.GET(path, svr.Handler())
	return engine
}

func Post(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmPost, handler)
}

func Delete(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmDelete, handler)
}

func Put(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmPut, handler)
}

func Head(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmHead, handler)
}

func Get(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmGet, handler)
}

func Options(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmOptions, handler)
}

func GetPost(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmGetPost, handler)
}

func All(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmAll, handler)
}

func PostWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmPost, header, versionName, handler)
}

func DeleteWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmDelete, header, versionName, handler)
}

func PutWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmPut, header, versionName, handler)
}

func HeadWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmHead, header, versionName, handler)
}

func GetWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmGet, header, versionName, handler)
}

func OptionsWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmOptions, header, versionName, handler)
}

func GetPostWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmGetPost, header, versionName, handler)
}

func AllWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmAll, header, versionName, handler)
}

func Use(middleware ...gin.HandlerFunc) {
	if engine != nil {
		engine.Use(middleware...)
	}
}

func getPathAppendApiModel(path string) string {
	// 获取 api-module
	apiModel := isc.ISCString(config.GetValueString("api-module")).Trim("/")
	// 获取api前缀
	ap := isc.ISCString(config.GetValueStringDefault("base.api.prefix", "")).Trim("/")
	if ap != "" {
		ApiPrefix = "/" + string(ap)
	}
	p2 := isc.ISCString(path).Trim("/")
	if strings.HasPrefix(string(p2), "api") {
		return fmt.Sprintf("/%s", p2)
	} else {
		return fmt.Sprintf("/%s/%s/%s", ApiPrefix, apiModel, p2)
	}
}
