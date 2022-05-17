package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/cron"
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
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("准备好了吗，我要panic了")
			//panic("panic了,打我吗？哈哈哈哈")
			go func(idx int) {
				c_s := cron.New()
				_ = c_s.AddFunc("*/1 * * * * ?", func() {
					panic("panic了,打我吗？哈哈哈哈")

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
		c.Data(200, "text/plain", []byte("hello"))
	})
	server.StartServer()
}
