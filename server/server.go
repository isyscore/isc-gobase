package server

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/server/rsp"
	"sync"

	//"github.com/isyscore/isc-gobase/tracing"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/listener"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

var GoBaseVersion = "1.4.3"
var ApiPrefix = "/api"

var engine *gin.Engine = nil
var pprofHave = false
var requestStorage goid.LocalStorage

var loadLock sync.Mutex
var serverLoaded = false

//type methodTrees []methodTree

var ginHandlers []gin.HandlerFunc

func init() {
	isc.PrintBanner()
	config.LoadConfig()
	printVersionAndProfile()
	requestStorage = goid.NewLocalStorage()
}

// 提供给外部注册使用
func AddGinHandlers(handler gin.HandlerFunc) {
	if nil == ginHandlers {
		var ginHandlersTem []gin.HandlerFunc
		ginHandlers = ginHandlersTem
	}

	ginHandlers = append(ginHandlers, handler)
}

func InitServer() {
	loadLock.Lock()
	defer loadLock.Unlock()
	if serverLoaded {
		return
	}
	if !config.ExistConfigFile() || !config.GetValueBoolDefault("base.server.enable", false) {
		return
	}

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
		gin.DefaultWriter = io.Discard
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	engine = gin.New()

	if config.GetValueBoolDefault("base.debug.enable", true) {
		// 注册pprof
		if config.GetValueBoolDefault("base.server.gin.pprof.enable", false) {
			pprofHave = true
			pprof.Register(engine)
		}
	}
	engine.Use(Cors(), gin.Recovery())
	engine.Use(RequestSaveHandler())
	engine.Use(rsp.ResponseHandler())

	for _, handler := range ginHandlers {
		engine.Use(handler)
	}

	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("base.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(apiPreAndModule())
	}

	if config.GetValueBoolDefault("base.debug.enable", true) {
		// 注册 配置查看和变更功能
		if config.GetValueBoolDefault("base.endpoint.config.enable", false) {
			RegisterConfigWatchEndpoint(apiPreAndModule())
		}

		// 注册 bean管理的功能
		if config.GetValueBoolDefault("base.endpoint.bean.enable", false) {
			RegisterBeanWatchEndpoint(apiPreAndModule())
		}

		// 注册 debug的帮助命令
		RegisterHelpEndpoint(apiPreAndModule())
	}

	// 注册 swagger的功能
	if config.GetValueBoolDefault("base.swagger.enable", false) {
		RegisterSwaggerEndpoint()
	}

	// 添加配置变更事件的监听
	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)

	appName := config.GetValueStringDefault("base.application.name", "isc-gobase")

	logger.InitLog(appName)
	serverLoaded = true
}

func ConfigChangeListener(event listener.BaseEvent) {
	ev := event.(listener.ConfigChangeEvent)
	if ev.Key == "base.server.gin.pprof.enable" {
		if isc.ToBool(ev.Value) && !pprofHave {
			pprofHave = true
			pprof.Register(engine)
		}
	}
}

func apiPreAndModule() string {
	ap := config.GetValueStringDefault("base.api.prefix", "")
	if ap != "" {
		ApiPrefix = ap
	}
	return ApiPrefix + "/" + config.ApiModule
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

	if engine == nil {
		return
	}

	listener.PublishEvent(listener.ServerRunStartEvent{})

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
	listener.PublishEvent(listener.ServerRunFinishEvent{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("服务端准备关闭...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := engineServer.Shutdown(ctx); err != nil {
		logger.Warn("服务关闭异常: %v", err.Error())
	}
	logger.Warn("服务端退出")
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

func Engine() *gin.Engine {
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
	RegisterRoute(apiBase+"/config/values/yaml", HmGet, config.GetConfigDeepValues)
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

func RegisterSwaggerEndpoint() gin.IRoutes {
	RegisterRoute("/swagger/*any", HmGet, ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine
}

func RegisterHelpEndpoint(apiBase string) gin.IRoutes {
	if "" == apiBase {
		return nil
	}
	RegisterRoute(apiBase+"/debug/help", HmGet, debug.Help)
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
		InitServer()
		return true
	}
	return true
}

func RegisterRoute(path string, method HttpMethod, handler gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	if engine == nil {
		logger.Warn("server没有启动，请配置 base.server.enable 或者查看相关日志")
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

func RequestSaveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestStorage.Set(c.Request)
		logger.PutHead(c.Request.Header)
	}
}

func GetRequest() *http.Request {
	return requestStorage.Get().(*http.Request)
}

func GetHeader() http.Header {
	req := requestStorage.Get()
	if req == nil {
		return nil
	}
	reqS := req.(*http.Request)
	return reqS.Header
}

func GetRemoteAddr() string {
	req := requestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.RemoteAddr
}

func GetHeaderWithKey(headKey string) string {
	req := requestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.Header.Get(headKey)
}
