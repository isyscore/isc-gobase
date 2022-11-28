package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/cron"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/server/rsp"
	"github.com/isyscore/isc-gobase/server/test/pojo"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server"
)

func TestServer(t *testing.T) {
	server.RegisterCustomHealthCheck("/api/sample",
		func() string {
			return "OK"
		},
		func() string {
			return "OK"
		},
		func() string {
			return "OK"
		},
	)

	logger.Info("server started")
	//panic("和水水水水")
	////log.Fatal("")
	go func() {
		for i := 0; i < 100; i++ {
			//panic("打我吗？")
			go func(idx int) {
				c_s := cron.New()
				_ = c_s.AddFunc("*/1 * * * * ?", func() {
					_, _ = fmt.Fprintf(os.Stderr, "我好帅，%s\n", "哈哈哈")
					_, _ = fmt.Fprintf(os.Stdout, "是真的\n")
					//logger.Debug("协程ID=：%d,我是库陈胜Debug", idx)
					logger.Info("协程ID=：%d,我是库陈胜Info", idx)
					//logger.Warn("协程ID=：%d,我是库陈胜Warn", idx)
					//logger.Error("协程ID=：%d,我是库陈胜Error", idx)
					//logger.Panic("我可以写入了吗?")
					//logger.Fatal("我是fatal")
					//panic("打我吗？")
				})
				c_s.Start()
			}(i)
		}
	}()

	server.StartServer()
}

func TestApiVersion(t *testing.T) {
	fmt.Printf("step 1\n")
	server.RegisterRouteWithHeaders("/api/sample", server.HmGet, []string{"isc-api-version"}, []string{"1.0"}, func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("hello 1.0"))
	})
	fmt.Printf("step 2\n")
	server.RegisterRouteWithHeaders("/api/sample", server.HmGet, []string{"isc-api-version"}, []string{"2.0"}, func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("hello 2.0"))
	})
	fmt.Printf("step 3\n")
	server.RegisterRouteWithHeaders("/api/sample", server.HmGet, []string{"isc-api-version"}, []string{"3.0"}, func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("hello 3.0"))
	})
	server.StartServer()
}

func TestErrorPrint(t *testing.T) {
	server.RegisterRoute("/api/data", server.HmGet, func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("hello 3.0"))
	})
	server.StartServer()
}

func TestServerGet(t *testing.T) {
	server.Get("/info", func(c *gin.Context) {
		logger.Debug("debug的日志")
		logger.Info("info的日志")
		logger.Warn("warn的日志")
		logger.Error("error的日志")
		c.Data(200, "text/plain", []byte("hello"))
	})

	// 测试事件监听机制
	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.BaseEvent) {
		logger.Info("应用启动完成")
	})

	server.StartServer()
}

func TestServer2(t *testing.T) {
	server.Get("/test/req1", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("hello"))
	})

	server.Get("/test/req2", func(c *gin.Context) {
		rsp.SuccessOfStandard(c, "value")
	})

	server.Get("/test/req3/:key", func(c *gin.Context) {
		rsp.SuccessOfStandard(c, c.Param("key"))
	})

	server.Post("/test/rsp1", func(c *gin.Context) {
		testReq := pojo.TestReq{}
		_ = isc.DataToObject(c.Request.Body, &testReq)
		rsp.SuccessOfStandard(c, testReq)
	})

	server.Get("/test/err", func(c *gin.Context) {
		rsp.FailedOfStandard(c, 500, "异常")
	})

	server.Run()
}

func init() {
	// 添加服务器启动完成事件监听
	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.BaseEvent) {
		logger.Info("应用启动完成")
	})

	// 添加服务器启动完成事件监听
	listener.AddListener(listener.EventOfServerStop, func(event listener.BaseEvent) {
		logger.Info("应用退出")
	})
}

func TestServerOnProfileIsPprof(t *testing.T) {
	server.Use(ApiVersionInterceptor())
	server.Get("data", func(c *gin.Context) {
		rsp.SuccessOfStandard(c, "data")
	})

	server.Run()
}

func ApiVersionInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
