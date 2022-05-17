package test

import (
	"fmt"
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
			panic("打我吗？")
			//go func(idx int) {
			//	c_s := cron.New()
			//	_ = c_s.AddFunc("*/1 * * * * ?", func() {
			//		//fmt.Fprintf(os.Stderr, "我好帅，%s\n", "哈哈哈")
			//		//fmt.Fprintf(os.Stdout, "是真的\n")
			//		//logger.Debug("协程ID=：%d,我是库陈胜Debug", idx)
			//		//logger.Info("协程ID=：%d,我是库陈胜Info", idx)
			//		//logger.Warn("协程ID=：%d,我是库陈胜Warn", idx)
			//		//logger.Error("协程ID=：%d,我是库陈胜Error", idx)
			//		//logger.Panic("我可以写入了吗?")
			//		//logger.Fatal("我是fatal")
			//		//log.Panic("是的，我可以写入了")
			//		panic("打我吗？")
			//	})
			//	c_s.Start()
			//}(i)
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
		c.Data(200, "text/plain", []byte("hello"))
	})
	server.StartServer()
}
