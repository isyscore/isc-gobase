package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/cron"

	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server"
)

func TestServer(t *testing.T) {
	server.InitServer()

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
			go func(idx int) {
				c_s := cron.New()
				_ = c_s.AddFunc("*/1 * * * * ?", func() {
					logger.Debug("协程ID=：%d,我是库陈胜Debug", idx)
					logger.Info("协程ID=：%d,我是库陈胜Info", idx)
					logger.Warn("协程ID=：%d,我是库陈胜Warn", idx)
					logger.Error("协程ID=：%d,我是库陈胜Error", idx)
				})
				c_s.Start()
			}(i)
		}
	}()

	server.StartServer()
}
